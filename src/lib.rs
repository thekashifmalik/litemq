use std::env;
use std::error::Error;
use std::collections::HashMap;
use tokio::sync::mpsc::Receiver;
use tokio::sync::mpsc::Sender;
use tokio::sync::mpsc::channel;
use tokio::sync::Mutex;
use std::sync::Arc;



use tonic::{transport, Request, Response, Status};
use log::{info, warn, debug, error};

use generated::lite_mq_server::LiteMq;
use generated::lite_mq_server::LiteMqServer;
use generated::Nothing;
use generated::QueueId;
use generated::QueueLength;
use generated::EnqueueRequest;
use generated::DequeueResponse;


/// This module is populated by the build script and contains the generated protobufs & gRPC code.
mod generated {
    tonic::include_proto!("_");
}

pub struct Server{
    queues: Mutex<HashMap<String, Arc<Mutex<Queue>>>>,
}

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

impl Server {
    pub fn new() -> Self {
        env_logger::Builder::from_env(env_logger::Env::default().filter_or("LOG_LEVEL", "info")).format_target(false).init();
        info!("LiteMQ");
    	debug!("debug logging enabled");
        Server{
            queues: Mutex::new(HashMap::new()),
        }
    }

    pub async fn serve(self) {
        let mut port = 42090;
        if let Ok(_port) = env::var("PORT") {
            match _port.parse::<u16>() {
                Ok(p) => {
                    port = p;
                }
                Err(e) => {
                    warn!("invalid PORT: {}", e);
                }
            }
        }

        let addr = match format!("0.0.0.0:{}", port).parse() {
            Ok(addr) => addr,
            Err(e) => {
                error!("{}", e);
                return
            }
        };

        info!("listening on {}", addr);

        match transport::Server::builder()
            .add_service(LiteMqServer::new(self))
            .serve(addr).await {
                Err(e) => {
                    let underlying_error = e.source().unwrap();
                    error!("{}: {}", e, underlying_error);
                    return
                }
                _ => ()
            };
        }
}

#[tonic::async_trait]
impl LiteMq for Server {

    async fn health(&self, _request: Request<Nothing>) -> Result<Response<Nothing>, Status> {
        log::info!("HEALTH");
        Ok(Response::new(Nothing::default()))
    }

    async fn enqueue(&self, request: Request<EnqueueRequest>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        let decoded = String::from_utf8(r.data.clone()).unwrap();
        log::info!("ENQUEUE {} \"{}\"", r.queue, decoded);

        let mut queues = self.queues.lock().await;
        let queue_mutex = match queues.get(&r.queue) {
            Some(q) => q.clone(),
            None => {
                let q = Arc::new(Mutex::new(Queue::new()));
                queues.insert(r.queue.clone(), q.clone());
                q
            }
        };
        // Release the server lock here.
        drop(queues);
        let mut queue = queue_mutex.lock().await;
        let length = queue.enqueue(r.data).await;
        Ok(Response::new(QueueLength{count: length}))
    }

    async fn dequeue(&self, request: Request<QueueId>) -> Result<Response<DequeueResponse>, Status> {
        let r = request.into_inner();
        info!("DEQUEUE {}", r.queue);

        let mut queues = self.queues.lock().await;
        let queue_mutex = match queues.get(&r.queue) {
            Some(q) => q.clone(),
            None => {
                queues.insert(r.queue.clone(), Arc::new(Mutex::new(Queue::new())));
                queues.get(&r.queue).unwrap().clone()
            }
        };
        // Release the server lock here.
        drop(queues);
        let mut queue = queue_mutex.lock().await;
        let mut rx = match queue.dequeue_or_receiver() {
            Ok(data) => {
                return Ok(Response::new(DequeueResponse{data: data}));
            }
            Err(rx) => rx,
        };
        // Release the queue lock here to avoid a deadlock
        drop(queue);
        let data = match rx.recv().await {
            Some(data) => data,
            None => {
                error!("no data received");
                return Err(Status::internal("no data received"));
            }
        };
        Ok(Response::new(DequeueResponse{data: data}))
    }

    async fn length(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        info!("LENGTH {}", r.queue);
        let queues = self.queues.lock().await;
        let count = match queues.get(&r.queue) {
            Some(q) => q.lock().await.length(),
            None => 0,
        };
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn purge(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        info!("LENGTH {}", r.queue);
        let mut queues = self.queues.lock().await;
        let count = match queues.remove(&r.queue) {
            Some(q) => q.lock().await.length(),
            None => 0,
        };
        Ok(Response::new(QueueLength{count: count}))
    }
}
