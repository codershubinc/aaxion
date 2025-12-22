use crate::models::FileInfo;
use axum::{
    extract::{Multipart, Query},
    http::StatusCode,
    response::{Html, IntoResponse, Response},
    Json,
};
use serde::Deserialize;
use std::path::{Path, PathBuf};
use tokio::fs::File;
use tokio::io::AsyncWriteExt;

const UPLOAD_DIR: &str = "/home/swap/aaxion/";

#[derive(Deserialize)]
pub struct ListFilesQuery {
    pub dir: Option<String>,
}

pub async fn list_files(Query(params): Query<ListFilesQuery>) -> Response {
    let requested_dir = params.dir.unwrap_or_else(|| UPLOAD_DIR.to_string());
    println!("📂 Listing files in directory: {}", requested_dir);

    let path = Path::new(&requested_dir);

    if !requested_dir.starts_with(UPLOAD_DIR) || requested_dir.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            "❌ Access forbidden: You can only list files inside the upload directory",
        )
            .into_response();
    }

    let mut entries = match tokio::fs::read_dir(&path).await {
        Ok(e) => e,
        Err(_) => return Json(vec![] as Vec<FileInfo>).into_response(),
    };

    let mut files: Vec<FileInfo> = Vec::new();

    while let Ok(Some(entry)) = entries.next_entry().await {
        let path = entry.path();
        let size = entry.metadata().await.map(|m| m.len()).unwrap_or(0);

        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if !name.starts_with(".") {
                files.push(FileInfo {
                    name: name.to_string(),
                    is_dir: path.is_dir(),
                    size,
                    path: requested_dir.clone(),
                    raw_path: path.to_string_lossy().to_string(),
                });
            }
        } else {
            println!("⚠️  Skipping invalid filename: {:?}", path);
        }
    }

    Json(files).into_response()
}

pub async fn upload_file(mut multipart: Multipart) -> impl IntoResponse {
    // We loop through ALL fields until we find the one we want
    while let Ok(Some(mut field)) = multipart.next_field().await {
        let name = field.name().unwrap_or("").to_string();

        if name == "file" {
            // FIX: Handle missing filename safely
            let file_name = match field.file_name() {
                Some(n) => n.to_string(),
                None => continue, // Skip if it's not a file
            };

            let safe_filename = Path::new(&file_name).file_name().unwrap().to_string_lossy();
            let file_path = PathBuf::from(UPLOAD_DIR).join(safe_filename.as_ref());

            println!("⬇️  Streaming start: {}", safe_filename);

            let mut file = match File::create(&file_path).await {
                Ok(f) => f,
                Err(e) => {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("❌ Server Error: {}", e),
                    )
                        .into_response()
                }
            };

            while let Ok(Some(chunk)) = field.chunk().await {
                if let Err(e) = file.write_all(&chunk).await {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("❌ Write Error: {}", e),
                    )
                        .into_response();
                }
            }

            println!("✅ Streaming complete: {}", safe_filename);
            return Html(format!("✅ Upload complete: {}", safe_filename)).into_response();
        }
    }

    // FIX: Provide a return value if the loop finishes without finding a file
    (StatusCode::BAD_REQUEST, "❌ No file field found in request").into_response()
}
