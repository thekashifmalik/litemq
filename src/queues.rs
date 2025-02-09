
use tokio::io::AsyncReadExt;
use tokio::io::AsyncWriteExt;
use tokio::sync::Mutex;
use tokio::sync::mpsc::Receiver;
use tokio::sync::mpsc::Sender;
use tokio::sync::mpsc::channel;

use log::debug;
use log::error;


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
}

impl Queue for PersistentQueue {
    async fn length(&self) -> i64 {
        match fs::read_to_string(&self.path).await {
            Ok(f) => f.lines().count() as i64,
            Err(_) => 0,
        }
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
        // Add a newline to the data.
        let mut payload = data;
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
        if self.length().await > 0 {
            let f: String = fs::read_to_string(&self.path).await.unwrap();
            let mut lines = f.lines().collect::<Vec<&str>>();
            let data = lines.remove(0).as_bytes().to_vec();
            let mut out = lines.iter()
                .map(|n| n.to_string())
                .collect::<Vec<_>>()
                .join("\n");
            // Append newline before write if needed.
            if out.len() > 0 && out.chars().last().unwrap() != '\n' {
                out.push('\n');
            }
            fs::write(&self.path, out).await.unwrap();
            return Some(data);
        } else {
            None
        }
    }

    async fn purge(&mut self) -> i64 {
        let length: i64 = self.length().await;
        match fs::remove_file(&self.path).await {
            Ok(_) => {},
            Err(e) => {
                error!("could not remove file {}: {}", self.path, e);
            }
        };
        return length;
    }
}
