use litemq::Server;


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let server = Server::new();
    server.serve().await;
    Ok(())
}
