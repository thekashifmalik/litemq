
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
    pub fn new(path: &str) -> Self {
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
        let f = match fs::read_to_string(&self.path).await {
            Ok(f) => f,
            Err(_) => "".to_string(),
        };
        let mut lines = f.lines().collect::<Vec<&str>>();
        let new = String::from_utf8(data).unwrap();
        lines.push(&new);
        let out = lines.iter()
            .map(|n| n.to_string())
            .collect::<Vec<_>>()
            .join("\n");
        match fs::write(&self.path, out.clone()).await {
            Ok(_) => {},
            Err(_) => {
                let mut file = File::create(&self.path).await.unwrap();
                file.write_all(out.as_bytes()).await.unwrap();
            }
        };
        self.length().await
    }

    async fn dequeue(&mut self) -> Option<Vec<u8>> {
        if self.length().await > 0 {
            let f: String = fs::read_to_string(&self.path).await.unwrap();
            let mut lines = f.lines().collect::<Vec<&str>>();
            let data = lines.remove(0).as_bytes().to_vec();
            let out = lines.iter()
                .map(|n| n.to_string())
                .collect::<Vec<_>>()
                .join("\n");
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
