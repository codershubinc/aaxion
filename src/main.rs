use aaxion::db;
use aaxion::routes::create_router;
use aaxion::services::discovery_service;
use axum::Extension;
use std::error::Error;
use std::net::SocketAddr;
use tokio::fs;

const UPLOAD_DIR: &str = "/home/swap/aaxion/";
const DB_PATH: &str = "/home/swap/aaxion/aaxion.db";
const PORT: u16 = 18875; // Used for both TCP (Web) and UDP (Discovery)

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    // 1. Ensure upload dir exists...
    if fs::metadata(UPLOAD_DIR).await.is_err() {
        fs::create_dir_all(UPLOAD_DIR).await?;
    }

    // Ensure database file exists...
    if fs::metadata(DB_PATH).await.is_err() {
        println!(
            "🆕 Database not found, creating new database at {}",
            DB_PATH
        );
        fs::File::create(DB_PATH).await?;
    }

    // Initialize DB pool (propagate errors)
    let pool = db::initialize_db_pool(DB_PATH.to_string()).await?;
    println!("✅ Database initialized at {}", DB_PATH);

    // Add default auth token (propagate errors)
    // println!("Add default auth token to database...");
    // db::add_auth_token(&pool).await?; run Once to add it

    // 2. Start Discovery on the SAME PORT
    tokio::spawn(async move {
        discovery_service::start_discovery_listener(PORT).await;
    });

    // 3. Start Web Server and attach DB pool as an Extension so handlers can use it
    let app = create_router().layer(Extension(pool.clone()));
    let addr = SocketAddr::from(([0, 0, 0, 0], PORT));

    println!("-------------------------------------------------");
    println!("🚀 Aaxion Server Running");
    println!("👉 Unified Port: {}", PORT);
    println!("-------------------------------------------------");

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}
