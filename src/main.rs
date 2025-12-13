use axum::{
    extract::{Multipart, DefaultBodyLimit},
    response::{Html, IntoResponse},
    routing::post,
    Router,
};
use std::net::SocketAddr;
use std::path::PathBuf;
use tokio::fs::{self, File};
use tokio::io::AsyncWriteExt; // Traits for writing to files
use tower_http::services::{ServeDir, ServeFile};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

// Configuration
const UPLOAD_DIR: &str = "./uploads";
const PORT: u16 = 8080;

#[tokio::main]
async fn main() {
    // 1. Initialize Logging
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new("localdrive_rs=debug,tower_http=debug"))
        .with(tracing_subscriber::fmt::layer())
        .init();

    // 2. Create Upload Directory
    if fs::metadata(UPLOAD_DIR).await.is_err() {
        fs::create_dir(UPLOAD_DIR).await.expect("Failed to create upload dir");
    }

    // 3. Define Routes
    // Note: We increase the body limit to unlimited (or very high) for large ISOs
    let app = Router::new()
        // Serve the UI from the static file
        .route_service("/", ServeFile::new("assets/index.html"))
        // Serve static assets (css, js, etc.)
        .nest_service("/assets", ServeDir::new("assets"))
        // Serve the static files (Browse view) from the current directory
        .nest_service("/files", ServeDir::new("."))
        // The Upload Logic
        .route("/upload", post(upload_handler))
        .layer(DefaultBodyLimit::disable()); // Disable default 2MB limit for streaming

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

// ---------------------------------------------------------
// HANDLERS
// ---------------------------------------------------------



// Handle Streaming Uploads
async fn upload_handler(mut multipart: Multipart) -> impl IntoResponse {
    // Iterate over the fields in the multipart form
    while let Ok(Some(mut field)) = multipart.next_field().await {
        let name = field.name().unwrap().to_string();
        
        // We only care about the "file" field
        if name == "file" {
            let file_name = field.file_name().unwrap().to_string();
            println!("⬇️  Streaming start: {}", file_name);

            let file_path = PathBuf::from(UPLOAD_DIR).join(&file_name);
            
            // Create the file on disk
            let mut file = match File::create(&file_path).await {
                Ok(f) => f,
                Err(e) => return Html(format!("❌ Server Error: {}", e)),
            };

            // STREAMING COPY:
            // Read chunks from the network stream and write directly to disk
            while let Ok(Some(chunk)) = field.chunk().await {
                if let Err(e) = file.write_all(&chunk).await {
                    return Html(format!("❌ Write Error: {}", e));
                }
            }
            
            println!("✅ Streaming complete: {}", file_name);
        }
    }

    Html("✅ Upload Complete".to_string())
}

