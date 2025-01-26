use tokio::sync::mpsc::Receiver;
use tokio::sync::mpsc::Sender;
use tokio::sync::mpsc::channel;

use log::debug;


pub struct Queue{
    pub messages: Vec<Vec<u8>>,
    pub channels: Vec<Sender<Vec<u8>>>,
}

impl Queue {
    pub fn new() -> Self {
        Queue{
            messages: Vec::new(),
            channels: Vec::new(),
        }
    }
}
impl Queue {
    pub fn length(&self) -> i64 {
        self.messages.len() as i64
    }

    pub async fn enqueue(&mut self, data: Vec<u8>) -> i64 {
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
                return self.length();
            }
        }
        // If we did not find a valid channel, we need to put the data into the queue.
        debug!("> {} bytes", data.len());
        self.messages.push(data);
        self.length()
    }

    pub fn dequeue_or_receiver(&mut self) -> Result<Vec<u8>, Receiver<Vec<u8>>> {
        if self.length() > 0 {
            let data = self.messages.remove(0);
            debug!("< {} bytes", data.len());
            return Ok(data);
        }
        let (tx, rx) = channel(1);
        self.channels.push(tx);
        debug!("* waiting");
        Err(rx)
    }
}
