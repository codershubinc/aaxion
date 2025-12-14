use axum::{
    extract::{Request, State},
    response::IntoResponse,
    routing::{any, get, post},
    Router,
};
use dav_server::{fakels::FakeLs, localfs::LocalFs, DavHandler};
use std::{net::SocketAddr, sync::Arc};
use tokio::fs;
use tower_http::services::{ServeDir, ServeFile};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod handlers;

const UPLOAD_DIR: &str = "./uploads";
const PORT: u16 = 8080;

#[tokio::main]
async fn main() {
    // 1. Logging
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            "localdrive_rs=debug,tower_http=debug",
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    // 2. Create Directory
    if fs::metadata(UPLOAD_DIR).await.is_err() {
        fs::create_dir(UPLOAD_DIR)
            .await
            .expect("Failed to create upload dir");
    }

    // 3. Setup WebDAV Handler (The "Drive" Logic)
    let webdav = DavHandler::builder()
        .filesystem(LocalFs::new(UPLOAD_DIR, false, false, false)) // Map ./uploads
        .locksystem(FakeLs::new())
        .strip_prefix("/webdav") // Important: Tell it we are serving under /webdav
        .build_handler();

    // Wrap it in Arc so we can share it across threads
    let dav_server = Arc::new(webdav);

    // 4. Define Routes
    let app = Router::new()
        // --- EXISTING WEB UI ROUTES ---
        .route_service("/", ServeFile::new("assets/index.html"))
        .route_service("/files", ServeFile::new("assets/files.html"))
        .nest_service("/assets", ServeDir::new("assets"))
        .route("/api/files/list", get(handlers::list_files))
        .route(
            "/api/files/download/:filename",
            get(handlers::download_file),
        )
        .route("/upload", post(handlers::upload_handler))
        // --- NEW WEBDAV ROUTE ---
        // We use 'any' because WebDAV uses method like PROPFIND, MKCOL, etc.
        .route("/webdav", any(webdav_handler))
        .route("/webdav/*path", any(webdav_handler))
        .with_state(dav_server);

    // 5. Start Server
    let addr = SocketAddr::from(([0, 0, 0, 0], PORT));
    println!("-------------------------------------------------");
    println!("🚀 Aaxion Server Running");
    println!("👉 Web UI:      http://{}:{}", "YOUR_IP", PORT);
    println!("👉 Drive Path:  http://{}:{}/webdav", "YOUR_IP", PORT);
    println!("⚠️  WINDOWS USERS: If you get 'Network name cannot be found':");
    println!("   1. Ensure WebClient service is running (services.msc)");
    println!("   2. Set HKLM\\SYSTEM\\CurrentControlSet\\Services\\WebClient\\Parameters\\BasicAuthLevel to 2");
    println!("-------------------------------------------------");

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// The "Glue" Function: Passes Axum requests to the WebDAV engine
async fn webdav_handler(State(dav): State<Arc<DavHandler>>, req: Request) -> impl IntoResponse {
    dav.handle(req).await
}
