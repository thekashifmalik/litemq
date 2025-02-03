use std::env;
use std::error::Error;
use std::collections::HashMap;
use tokio::fs;
use tokio::sync::Mutex;
use tokio::sync::mpsc::Receiver;
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
use queues::Queue;
use queues::InMemoryQueue;

mod queues;

/// This module is populated by the build script and contains the generated protobufs & gRPC code.
mod generated {
    tonic::include_proto!("_");
}

pub struct Server{
    queues: Mutex<HashMap<String, Arc<Mutex<InMemoryQueue>>>>,
}

impl Server {
    pub async fn new() -> Self {
        env_logger::Builder::from_env(env_logger::Env::default().filter_or("LOG_LEVEL", "info")).format_target(false).init();
        info!("LiteMQ");
        debug!("debug logging enabled");

        let data_dir =  match env::var("LITEMQ_DATA_DIR") {
            Ok(d) => {
                info!("LITEMQ_DATA_DIR: {}", d);
                fs::create_dir_all(&d).await.unwrap();
                Some(d)

            },
            Err(_) => {
                warn!("LITEMQ_DATA_DIR not set, using in-memory queues");
                None
            }
        };
        Self{
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

    async fn lock_and_get_or_create_queue(&self, name: &str) -> Arc<Mutex<InMemoryQueue>> {
        let mut locked = self.queues.lock().await;
        if let Some(q) = locked.get(name) {
            return q.clone();
        }
        // We have to create a new queue.
        let queue = Arc::new(Mutex::new(InMemoryQueue::new()));
        locked.insert(name.to_string(), queue.clone());
        queue
    }

    async fn lock_and_enqueue(&self, name: &str, data: Vec<u8>) -> i64 {
        let queue = self.lock_and_get_or_create_queue(name).await;
        let mut locked = queue.lock().await;
        locked.enqueue(data).await
    }

    async fn lock_and_dequeue_or_receiver(&self, name: &str) -> Result<Vec<u8>, Receiver<Vec<u8>>> {
        let queue = self.lock_and_get_or_create_queue(name).await;
        let mut locked = queue.lock().await;
        locked.dequeue_or_receiver().await
    }

    async fn lock_and_get_queue_length(&self, name: &str) -> i64 {
        let locked_1 = self.queues.lock().await;
        match locked_1.get(name) {
            None => 0,
            Some(queue) => {
                let locked_2 = queue.lock().await;
                locked_2.length().await
            },
        }
    }

    async fn lock_and_purge_queue(&self, name: &str) -> i64 {
        let locked_1 = self.queues.lock().await;
        match locked_1.get(name) {
            None => 0,
            Some(queue) => {
                let mut locked_2 = queue.lock().await;
                locked_2.purge().await
            },
        }
    }

    async fn lock_and_flush(&self) {
        let mut locked = self.queues.lock().await;
        locked.clear();
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
        let decoded: String = String::from_utf8(r.data.clone()).unwrap();
        log::info!("ENQUEUE {} \"{}\"", r.queue, decoded);

        let length = self.lock_and_enqueue(&r.queue, r.data).await;
        Ok(Response::new(QueueLength{count: length}))
    }

    async fn dequeue(&self, request: Request<QueueId>) -> Result<Response<DequeueResponse>, Status> {
        let r = request.into_inner();
        info!("DEQUEUE {}", r.queue);

        let mut rx = match self.lock_and_dequeue_or_receiver(&r.queue).await {
            Ok(data) => {
                return Ok(Response::new(DequeueResponse{data: data}));
            }
            Err(rx) => {
                rx
            },
        };
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
        let count = self.lock_and_get_queue_length(&r.queue).await;
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn purge(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        info!("PURGE {}", r.queue);
        let count = self.lock_and_purge_queue(&r.queue).await;
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn flush(&self, _: Request<Nothing>) -> Result<Response<Nothing>, Status> {
        info!("FLUSH");
        self.lock_and_flush().await;
        Ok(Response::new(Nothing::default()))
    }
}
