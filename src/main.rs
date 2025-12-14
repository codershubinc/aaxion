use axum::{
    extract::DefaultBodyLimit,
    routing::{get, post},
    Router,
};
use std::net::SocketAddr;
use tokio::fs;
use tower_http::services::{ServeDir, ServeFile};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

// Register the module
mod handlers;

// Configuration
const UPLOAD_DIR: &str = "./uploads";
const PORT: u16 = 8080;

#[tokio::main]
async fn main() {
    // 1. Initialize Logging
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            "localdrive_rs=debug,tower_http=debug",
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    // 2. Create Upload Directory
    if fs::metadata(UPLOAD_DIR).await.is_err() {
        fs::create_dir(UPLOAD_DIR)
            .await
            .expect("Failed to create upload dir");
    }

    // 3. Define Routes
    let app = Router::new()
        .route_service("/", ServeFile::new("assets/index.html"))
        .route_service("/files", ServeFile::new("assets/files.html"))
        .nest_service("/assets", ServeDir::new("assets"))
        .nest_service("/raw", ServeDir::new(UPLOAD_DIR))
        // USE THE MODULE HERE: handlers::function_name
        .route("/api/files/list", get(handlers::list_files))
        .route(
            "/api/files/download/:filename",
            get(handlers::download_file),
        )
        .route("/upload", post(handlers::upload_handler))
        .layer(DefaultBodyLimit::disable());

    // 4. Start Server
    let addr = SocketAddr::from(([0, 0, 0, 0], PORT));
    println!("-------------------------------------------------");
    println!("🚀 LocalDrive-RS Running");
    println!("👉 UI & Upload: http://{}:{}", "YOUR_LAN_IP", PORT);
    println!("👉 File Browser: http://{}:{}/files", "YOUR_LAN_IP", PORT);
    println!("-------------------------------------------------");

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
