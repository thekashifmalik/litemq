use std::error::Error;
use std::sync::RwLock;
use std::collections::HashMap;
use std::sync::mpsc::Sender;
use std::sync::mpsc::Receiver;



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
    queues: RwLock<HashMap<String, Queue>>,
}

pub struct Queue{
    pub messages: Vec<Vec<u8>>,
}

impl Queue {
    pub fn new() -> Self {
        Queue{
            messages: Vec::new(),
        }
    }
}

impl Server {
    pub fn new() -> Self {
        env_logger::Builder::from_env(env_logger::Env::default().default_filter_or("info")).init();
        info!("LiteMQ");
    	debug!("debug logging enabled");
        Server{
            queues: RwLock::new(HashMap::new()),
        }
    }

    pub async fn serve(self) {
        let addr = match "0.0.0.0:42069".parse() {
            Ok(addr) => addr,
            Err(e) => {
                error!("{}", e);
                return
            }
        };

        info!("binding to {}", addr);

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

        let mut queues = match self.queues.write() {
            Ok(queues) => queues,
            Err(e) => {
                warn!("lock poisoned: {}", e);
                e.into_inner()
            }
        };
        let queue = match queues.get_mut(&r.queue) {
            Some(q) => q,
            None => {
                queues.insert(r.queue.clone(), Queue::new());
                queues.get_mut(&r.queue).unwrap()
            }
        };
        queue.messages.push(r.data);
        let count = queue.messages.len() as i64;
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn dequeue(&self, request: Request<QueueId>) -> Result<Response<DequeueResponse>, Status> {
        let r = request.into_inner();
        info!("DEQUEUE {}", r.queue);

        let mut queues = match self.queues.write() {
            Ok(queues) => queues,
            Err(e) => {
                warn!("lock poisoned: {}", e);
                e.into_inner()
            }
        };
        let queue = match queues.get_mut(&r.queue) {
            Some(q) => q,
            None => {
                queues.insert(r.queue.clone(), Queue::new());
                queues.get_mut(&r.queue).unwrap()
            }
        };
        // TODO: Implement blocking reads using channels here.
        let data = queue.messages.remove(0);
        Ok(Response::new(DequeueResponse{data: data}))
    }

    async fn length(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        info!("LENGTH {}", r.queue);
        let queues = match self.queues.read() {
            Ok(queues) => queues,
            Err(e) => {
                warn!("lock poisoned: {}", e);
                e.into_inner()
            }
        };
        let count = match queues.get(&r.queue) {
            Some(q) => q.messages.len() as i64,
            None => 0,
        };
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn purge(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        info!("LENGTH {}", r.queue);
        let mut queues = match self.queues.write() {
            Ok(queues) => queues,
            Err(e) => {
                warn!("lock poisoned: {}", e);
                e.into_inner()
            }
        };
        let count = match queues.remove(&r.queue) {
            Some(q) => q.messages.len() as i64,
            None => 0,
        };
        Ok(Response::new(QueueLength{count: count}))
    }
}
