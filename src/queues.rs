
use tokio::io::AsyncReadExt;
use tokio::io::AsyncWriteExt;
use tokio::sync::Mutex;
use tokio::sync::mpsc::Receiver;
use tokio::sync::mpsc::Sender;
use tokio::sync::mpsc::channel;

use log::debug;


pub trait Queue {
    async fn length(&self) -> i64;
    async fn enqueue(&mut self, data: Vec<u8>) -> i64;
    async fn dequeue_or_receiver(&mut self) -> Result<Vec<u8>, Receiver<Vec<u8>>>;
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
        if self.channels.len() > 0 {
            let mut tx = self.channels.remove(0);
            // Currently dequeue is not cleaning up senders that are closed so we need to do it here.
            while tx.is_closed() && self.channels.len() > 0 {
                tx = self.channels.remove(0);
            }
            // If we found a valid channel, we can send the data and return the queue length.
            if !tx.is_closed() {
                debug!("* {} bytes", data.len());
                tx.send(data).await.unwrap();
                return self.length().await;
            }
        }
        // If we did not find a valid channel, we need to put the data into the queue.
        debug!("> {} bytes", data.len());
        self.messages.push(data);
        self.length().await
    }

    async fn dequeue_or_receiver(&mut self) -> Result<Vec<u8>, Receiver<Vec<u8>>> {
        if self.length().await > 0 {
            let data = self.messages.remove(0);
            debug!("< {} bytes", data.len());
            return Ok(data);
        }
        let (tx, rx) = channel(1);
        self.channels.push(tx);
        debug!("* waiting");
        Err(rx)
    }

    async fn purge(&mut self) -> i64 {
        let length: i64 = self.length().await;
        self.messages.clear();
        return length;
    }
}
