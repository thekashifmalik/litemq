use std::env;
use std::error::Error;
use std::collections::HashMap;
use tokio::fs;
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
use queues::ThreadSafeQueue;
use queues::InMemoryQueue;

mod queues;

/// This module is populated by the build script and contains the generated protobufs & gRPC code.
mod generated {
    tonic::include_proto!("_");
}

pub struct Server{
    queues: Mutex<HashMap<String, Arc<ThreadSafeQueue>>>,
    // data_dir: Option<String>,
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
        Server{
            queues: Mutex::new(HashMap::new()),
            // data_dir: data_dir,
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

    async fn lock_and_get_or_create_queue(&self, name: &str) -> Arc<ThreadSafeQueue> {
        let mut queues = self.queues.lock().await;
        match queues.get(name) {
            Some(q) => q.clone(),
            None => {
                let q = Arc::new(ThreadSafeQueue::new(InMemoryQueue::new()));
                queues.insert(name.to_string(), q.clone());
                q
            }
        }
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

        let queue = self.lock_and_get_or_create_queue(&r.queue).await;
        let length = queue.lock_and_enqueue(r.data).await;
        Ok(Response::new(QueueLength{count: length}))
    }

    async fn dequeue(&self, request: Request<QueueId>) -> Result<Response<DequeueResponse>, Status> {
        let r = request.into_inner();
        info!("DEQUEUE {}", r.queue);

        let queue = self.lock_and_get_or_create_queue(&r.queue).await;
        let mut rx = match queue.lock_and_dequeue_or_receiver().await {
            Ok(data) => {
                return Ok(Response::new(DequeueResponse{data: data}));
            }
            Err(rx) => rx,
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
        let queue = self.lock_and_get_or_create_queue(&r.queue).await;
        let count = queue.lock_and_length().await;
        Ok(Response::new(QueueLength{count: count}))
    }

    async fn purge(&self, request: Request<QueueId>) -> Result<Response<QueueLength>, Status> {
        let r = request.into_inner();
        info!("LENGTH {}", r.queue);
        let mut queues = self.queues.lock().await;
        let count = match queues.remove(&r.queue) {
            Some(q) => q.lock_and_length().await,
            None => 0,
        };
        Ok(Response::new(QueueLength{count: count}))
    }
}
