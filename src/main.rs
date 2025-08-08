use litemq::Server;
use std::error::Error;


#[tokio::main]
async fn main() -> Result<(), Box<dyn Error + Send + Sync>> {
    let server = Server::new().await?;
    server.serve().await;
    Ok(())
}
