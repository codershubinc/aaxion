use aaxion::routes::create_router;
use std::net::SocketAddr;
use tokio::fs; // Use tokio instead of futures for fs

const UPLOAD_DIR: &str = "/home/swap/aaxion/";
const PORT: u16 = 18875;

#[tokio::main]
async fn main() {
    // 1. Ensure the upload directory exists
    if fs::metadata(UPLOAD_DIR).await.is_err() {
        fs::create_dir_all(UPLOAD_DIR) // create_dir_all is safer for absolute paths
            .await
            .expect("Failed to create upload dir");
    }

    // 2. Initialize the router
    let app = create_router();

    // 3. Define the address
    let addr = SocketAddr::from(([0, 0, 0, 0], PORT));

    println!("-------------------------------------------------");
    println!("🚀 Aaxion Server Running");
    println!("👉 Web UI:      http://localhost:{}", PORT);
    println!("👉 Upload Dir:  {}", UPLOAD_DIR);
    println!("-------------------------------------------------");

    // 4. START THE SERVER (This was missing)
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
