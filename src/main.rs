use litemq::Server;
use std::error::Error;


#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let server = Server::new();
    server.serve().await;
    Ok(())
}
