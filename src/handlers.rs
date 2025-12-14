use axum::{
    extract::{Multipart, Path},
    response::{Html, IntoResponse, Json},
    http::{header, StatusCode},
    body::{Body, Bytes},
};
use serde::Serialize;
use std::path::PathBuf;
use tokio::fs::File;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use futures::stream;

const UPLOAD_DIR: &str = "./uploads";

#[derive(Serialize)]
pub struct FileInfo {
    pub name: String,
    pub is_dir: bool,
    pub size: u64,
    pub path: String,
}

// ---------------------------------------------------------
// HANDLERS
// ---------------------------------------------------------

pub async fn list_files() -> impl IntoResponse {
    let mut entries = match tokio::fs::read_dir(UPLOAD_DIR).await {
        Ok(e) => e,
        Err(_) => return Json(vec![]),
    };

    let mut files = Vec::new();

    while let Ok(Some(entry)) = entries.next_entry().await {
        let path = entry.path();

        let size = entry.metadata().await.map(|m| m.len()).unwrap_or(0);
        let file_path = UPLOAD_DIR;
        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if !name.starts_with('.') {
                files.push(FileInfo {
                    name: name.to_string(),
                    is_dir: path.is_dir(),
                    size,
                    path: file_path.to_string().replace(".", ""),
                });
            }
        }
    }

    Json(files)
}


pub async fn upload_handler(mut multipart: Multipart) -> impl IntoResponse {
    while let Ok(Some(mut field)) = multipart.next_field().await {
        let name = field.name().unwrap().to_string();

        if name == "file" {
            let file_name = field.file_name().unwrap().to_string();
            println!("⬇️  Streaming start: {}", file_name);

            let file_path = PathBuf::from(UPLOAD_DIR).join(&file_name);

            let mut file = match File::create(&file_path).await {
                Ok(f) => f,
                Err(e) => return Html(format!("❌ Server Error: {}", e)),
            };

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

pub async fn download_file(Path(filename): Path<String>) -> impl IntoResponse {
    let file_path = PathBuf::from(UPLOAD_DIR).join(&filename);

    let file = match File::open(&file_path).await {
        Ok(f) => f,
        Err(_) => return (StatusCode::NOT_FOUND, "File not found").into_response(),
    };

    let stream = stream::unfold(file, |mut file| async move {
        let mut buf = vec![0u8; 64 * 1024]; // 64KB chunk
        match file.read(&mut buf).await {
            Ok(0) => None,
            Ok(n) => {
                buf.truncate(n);
                Some((Ok::<_, std::io::Error>(Bytes::from(buf)), file))
            }
            Err(e) => Some((Err(e), file)),
        }
    });

    let body = Body::from_stream(stream);

    let headers = [
        (header::CONTENT_TYPE, "application/octet-stream".to_string()),
        (header::CONTENT_DISPOSITION, format!("attachment; filename=\"{}\"", filename)),
    ];

    (headers, body).into_response()
}