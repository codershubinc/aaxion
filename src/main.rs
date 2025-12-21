use aaxion::routes::create_router;
use aaxion::services::discovery_service;
use std::net::SocketAddr;
use tokio::fs;

const UPLOAD_DIR: &str = "/home/swap/aaxion/";
const PORT: u16 = 18875; // Used for both TCP (Web) and UDP (Discovery)

#[tokio::main]
async fn main() {
    // 1. Ensure upload dir exists...
    if fs::metadata(UPLOAD_DIR).await.is_err() {
        fs::create_dir_all(UPLOAD_DIR).await.unwrap();
    }

    // 2. Start Discovery on the SAME PORT
    tokio::spawn(async move {
        discovery_service::start_discovery_listener(PORT).await;
    });

    // 3. Start Web Server
    let app = create_router();
    let addr = SocketAddr::from(([0, 0, 0, 0], PORT));

    println!("-------------------------------------------------");
    println!("🚀 Aaxion Server Running");
    println!("👉 Unified Port: {}", PORT);
    println!("-------------------------------------------------");

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
