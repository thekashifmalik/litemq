use std::borrow::Borrow;
use std::error::Error;
use std::sync::{Mutex, RwLock};
use std::collections::HashMap;


use tonic::{transport, Request, Response, Status};
use log::{info, warn, debug, error};

use generated::lite_mq_server::LiteMq;
use generated::lite_mq_server::LiteMqServer;
use generated::Nothing;
use generated::QueueId;
use generated::QueueLength;
use generated::EnqueueRequest;
use generated::DequeueResponse;


mod generated {
    tonic::include_proto!("_"); // The string specified here must match the proto package name
}

pub struct Server{
    messages: RwLock<HashMap<String, Vec<Vec<u8>>>>
}

impl Server {
    pub fn new() -> Self {
        env_logger::Builder::from_env(env_logger::Env::default().default_filter_or("info")).init();
        info!("LiteMQ");
    	debug!("debug logging enabled");
        Server{
            messages: RwLock::new(HashMap::new()),
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

    async fn health(&self, request: Request<Nothing>) -> Result<Response<Nothing>, Status> {
        log::info!("HEALTH");
        Ok(Response::new(Nothing::default()))
    }

    async fn enqueue(&self, request: Request<EnqueueRequest>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        let decoded = String::from_utf8(r.data.clone()).unwrap();
        log::info!("ENQUEUE {} \"{}\"", r.queue, decoded);
        let count;
        match self.messages.write() {
            Ok(mut messages) => {
                let queue = match messages.get_mut(&r.queue) {
                    Some(q) => q,
                    None => {
                        messages.insert(r.queue.clone(), Vec::new());
                        messages.get_mut(&r.queue).unwrap()
                    }
                };
                queue.push(r.data);
                count = queue.len() as i64;
            }
            Err(e) => {
                error!("{}", e);
                return Err(Status::internal("failed to lock messages"));
            }
        }
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn dequeue(&self, request: Request<QueueId>) -> Result<Response<DequeueResponse>, Status> {
        info!("{:?}", request);
        Ok(Response::new(DequeueResponse::default()))
    }

    async fn length(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        log::info!("LENGTH {}", r.queue);
        let count = match self.messages.read() {
            Ok(messages) => {
                match messages.get(&r.queue) {
                    Some(q) => q.len() as i64,
                    None => 0,
                }
            }
            Err(e) => {
                // TODO: Maybe we sold fail the request here instead.
                warn!("{}", e);
                0
            }
        };
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn purge(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        info!("{:?}", request);
        Ok(Response::new(QueueLength::default()))
    }
}
