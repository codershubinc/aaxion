use axum::{
    extract::{Multipart, DefaultBodyLimit},
    response::{Html, IntoResponse},
    routing::{get, post},
    Router,
};
use futures::TryStreamExt; // Helpers for streaming
use std::net::SocketAddr;
use std::path::PathBuf;
use tokio::fs::{self, File};
use tokio::io::AsyncWriteExt; // Traits for writing to files
use tower_http::services::ServeDir;
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
        // Serve the static files (Browse view) from the current directory
        .nest_service("/files", ServeDir::new("."))
        // The UI and Upload Logic
        .route("/", get(show_ui))
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

// Serve the HTML UI
async fn show_ui() -> Html<&'static str> {
    Html(HTML_TEMPLATE)
}

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

// ---------------------------------------------------------
// UI TEMPLATE (Embedded)
// ---------------------------------------------------------
const HTML_TEMPLATE: &str = r#"
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LocalDrive-RS</title>
    <style>
        :root { --bg: #111; --card: #1a1a1a; --text: #eee; --accent: #ff4d00; /* Rust Orange */ }
        body { background: var(--bg); color: var(--text); font-family: sans-serif; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        .container { background: var(--card); padding: 2rem; border-radius: 8px; width: 90%; max-width: 450px; text-align: center; border: 1px solid #333; }
        h2 { color: var(--accent); }
        .custom-file-btn { display: block; padding: 20px; border: 2px dashed #444; color: #888; border-radius: 8px; cursor: pointer; margin-bottom: 20px; transition: 0.2s; }
        .custom-file-btn:hover { border-color: var(--accent); color: var(--accent); }
        .btn { background: var(--accent); color: white; border: none; padding: 12px; width: 100%; font-weight: bold; cursor: pointer; border-radius: 4px; }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; }
        #progress { width: 100%; background: #333; height: 10px; margin-top: 15px; border-radius: 5px; overflow: hidden; display: none; }
        #bar { height: 100%; background: var(--accent); width: 0%; }
        #stats { display: flex; justify-content: space-between; font-size: 0.8rem; color: #666; margin-top: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>🦀 LocalDrive-RS</h2>
        <input type="file" id="file" style="display:none" onchange="updateLabel(this)">
        <label for="file" class="custom-file-btn" id="label">📂 Select File</label>
        <button class="btn" onclick="upload()" id="btn">Start Upload</button>
        
        <div id="progress"><div id="bar"></div></div>
        <div id="stats"><span id="speed">0 MB/s</span><span id="pct">0%</span></div>
        <br>
        <a href="/files" style="color: #666; text-decoration: none;">Browse Files &rarr;</a>
    </div>

    <script>
        function updateLabel(el) { if(el.files.length > 0) document.getElementById('label').innerText = el.files[0].name; }
        
        function upload() {
            let file = document.getElementById('file').files[0];
            if(!file) return alert("Pick a file!");
            
            let btn = document.getElementById('btn');
            btn.disabled = true;
            document.getElementById('progress').style.display = 'block';
            
            let xhr = new XMLHttpRequest();
            xhr.open("POST", "/upload", true);
            let start = Date.now();
            
            xhr.upload.onprogress = e => {
                if(e.lengthComputable) {
                    let pct = (e.loaded / e.total) * 100;
                    let sec = (Date.now() - start) / 1000;
                    let speed = sec > 0 ? (e.loaded/sec)/(1024*1024) : 0;
                    
                    document.getElementById('bar').style.width = pct + "%";
                    document.getElementById('pct').innerText = Math.round(pct) + "%";
                    document.getElementById('speed').innerText = speed.toFixed(2) + " MB/s";
                }
            };
            
            xhr.onload = () => { btn.disabled = false; alert("Done!"); };
            let fd = new FormData(); fd.append("file", file);
            xhr.send(fd);
        }
    </script>
</body>
</html>
"#;