
use std::os::unix::fs::MetadataExt;

use tokio::io::AsyncReadExt;
use tokio::io::AsyncWriteExt;
use tokio::sync::Mutex;
use tokio::sync::mpsc::Receiver;
use tokio::sync::mpsc::Sender;
use tokio::sync::mpsc::channel;

use log::debug;
use log::error;
use base64::{Engine as _, engine::general_purpose};


pub trait Queue {
    async fn length(&self) -> i64;
    async fn enqueue(&mut self, data: Vec<u8>) -> i64;
    async fn dequeue(&mut self) -> Option<Vec<u8>>;
    async fn purge(&mut self) -> i64;
}

pub struct InMemoryQueue{
    pub messages: Vec<Vec<u8>>,
    pub channels: Vec<Sender<Vec<u8>>>,
}

impl InMemoryQueue {
    pub fn new() -> Self {
        InMemoryQueue{
            messages: Vec::new(),
            channels: Vec::new(),
        }
    }
}

impl Queue for InMemoryQueue {
    async fn length(&self) -> i64 {
        self.messages.len() as i64
    }

    async fn enqueue(&mut self, data: Vec<u8>) -> i64 {
        self.messages.push(data);
        self.length().await
    }

    async fn dequeue(&mut self) -> Option<Vec<u8>> {
        if self.length().await > 0 {
            Some(self.messages.remove(0))
        } else {
            None
        }
    }

    async fn purge(&mut self) -> i64 {
        let length: i64 = self.length().await;
        self.messages.clear();
        return length;
    }
}


use tokio::fs::File;
use tokio::fs;


pub struct PersistentQueue{
    pub path: String,
    pub channels: Vec<Sender<Vec<u8>>>,
}

impl PersistentQueue {
    pub async fn new(path: &str) -> Self {
        fs::write(path, "").await.unwrap();
        PersistentQueue{
            path: path.to_string(),
            channels: Vec::new(),
        }
    }
    pub fn existing(path: &str) -> Self {
        PersistentQueue{
            path: path.to_string(),
            channels: Vec::new(),
        }
    }

    // Helper method to get the number of dequeued messages
    async fn get_num_dequeued(&self) -> u64 {
        let dequeue_path = self.get_dequeue_path();
        match fs::metadata(&dequeue_path).await {
            Ok(f) => f.size(),
            Err(_) => 0,
        }
    }

    // Helper method to get the dequeue tracking file path
    fn get_dequeue_path(&self) -> String {
        format!("{}.dequeued", self.path)
    }
}

impl Queue for PersistentQueue {
    async fn length(&self) -> i64 {
        let total_length = match fs::read_to_string(&self.path).await {
            Ok(f) => f.lines().count() as i64,
            Err(_) => 0,
        };

        // Check how many messages have already been dequeued
        let num_dequeued = self.get_num_dequeued().await as i64;

        // Return available messages (total - dequeued)
        total_length - num_dequeued
    }

    async fn enqueue(&mut self, data: Vec<u8>) -> i64 {
        // Get the current length of the queue so we can calculate the final length at the end of this function. We do
        // this because there is a race-condition between the file append/write and the length read.
        let length = self.length().await;
        // Open the file for appending
        let mut f = match fs::OpenOptions::new()
            .append(true)
            .open(&self.path)
            .await {
            Ok(f) => f,
            Err(e) => {
                error!("could not open file {}: {}", self.path, e);
                return 0;
            }
        };
        // Encode the data in base64 to avoid newline issues
        let encoded_data = general_purpose::STANDARD.encode(&data);
        // Add a newline to separate messages (safe since base64 never contains newlines)
        let mut payload = encoded_data.into_bytes();
        payload.push(b'\n');
        // Write the data to the file
        match f.write_all(&payload).await {
            Ok(_) => {},
            Err(e) => {
                error!("could not write to file {}: {}", self.path, e);
                return 0;
            }
        };
        // Return the new length of the queue. See start of method to see understand why we do this.
        length + 1
    }

    async fn dequeue(&mut self) -> Option<Vec<u8>> {
        // Get num_dequeued first
        let num_dequeued = self.get_num_dequeued().await;

        // Read the file once and compute both total length and get the line we need
        let f = match fs::read_to_string(&self.path).await {
            Ok(content) => content,
            Err(e) => {
                error!("could not read file {}: {}", self.path, e);
                return None;
            }
        };

        let lines: Vec<&str> = f.lines().collect();
        let total_length = lines.len() as i64;
        let available_length = total_length - num_dequeued as i64;

        if available_length > 0 {
            // Get the line at the current offset (num_dequeued position)
            if let Some(line) = lines.get(num_dequeued as usize) {
                // Decode the base64 encoded data back to original bytes
                let data = match general_purpose::STANDARD.decode(line) {
                    Ok(decoded) => decoded,
                    Err(e) => {
                        error!("could not decode base64 data from line {}: {}", num_dequeued, e);
                        return None;
                    }
                };

                // Update the dequeue counter by appending 'x'
                let dequeue_path = self.get_dequeue_path();
                match fs::OpenOptions::new()
                    .create(true)
                    .append(true)
                    .open(&dequeue_path)
                    .await {
                    Ok(mut f) => {
                        if let Err(e) = f.write_all(b"x").await {
                            error!("could not write to dequeue file {}: {}", dequeue_path, e);
                            return None;
                        }
                    },
                    Err(e) => {
                        error!("could not open dequeue file {}: {}", dequeue_path, e);
                        return None;
                    }
                };

                return Some(data);
            } else {
                error!("line at offset {} not found in file {}", num_dequeued, self.path);
                return None;
            }
        } else {
            // Queue is empty so we can compact to reclaim space
            if total_length > 0 && num_dequeued > 0 {
                debug!("compacting queue: {}", self.path);

                // Clear the main queue file
                match fs::write(&self.path, "").await {
                    Ok(_) => {
                        // Remove the dequeue tracking file if it exists
                        let dequeue_path = self.get_dequeue_path();
                        if fs::metadata(&dequeue_path).await.is_ok() {
                            if let Err(e) = fs::remove_file(&dequeue_path).await {
                                error!("failed to remove dequeue file {}: {}", dequeue_path, e);
                            }
                        }
                    },
                    Err(e) => {
                        error!("failed to compact queue {}: {}", self.path, e);
                    }
                }
            }
            None
        }
    }

    async fn purge(&mut self) -> i64 {
        let length: i64 = self.length().await;

        // Remove the main queue file
        match fs::remove_file(&self.path).await {
            Ok(_) => {},
            Err(e) => {
                error!("could not remove file {}: {}", self.path, e);
            }
        };

        // Also remove the dequeue tracking file
        let dequeue_path = self.get_dequeue_path();
        match fs::remove_file(&dequeue_path).await {
            Ok(_) => {},
            Err(_) => {
                // It's okay if the dequeue file doesn't exist
            }
        };

        return length;
    }
}
