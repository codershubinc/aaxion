use crate::models::FileInfo;
use axum::{
    body::Body,
    extract::{Multipart, Query},
    http::{HeaderMap, StatusCode},
    response::{IntoResponse, Response},
    Json,
};
use futures::StreamExt;
use serde::Deserialize;
use std::path::{Path, PathBuf};
use std::time::SystemTime;
use tokio::fs::File;
use tokio::io::{AsyncWriteExt, BufWriter};
use tokio_util::io::StreamReader;

/// Root upload directory
const UPLOAD_DIR: &str = "/home/swap/aaxion/";

#[derive(Deserialize)]
pub struct ListFilesQuery {
    pub dir: Option<String>,
    pub name: Option<String>, // optional original filename from client
}

#[derive(Deserialize)]
pub struct CreateFolderRequest {
    pub path: String,
}

#[derive(Deserialize)]
pub struct CreateFileRequest {
    pub path: String,
    pub content: Option<String>,
}

#[derive(Deserialize)]
pub struct DeleteRequest {
    pub path: String,
}

#[derive(Deserialize)]
pub struct DownloadQuery {
    pub path: String,
}

/// --- Helper utilities ---------------------------------------------------

fn sanitize_filename(raw: &str) -> String {
    // Keep only the basename to avoid directory traversal
    let base = Path::new(raw)
        .file_name()
        .and_then(|s| s.to_str())
        .unwrap_or("uploaded_file");

    // Allow only alphanumeric and a few safe punctuation chars; replace the rest with '_'
    let allowed = |c: char| c.is_ascii_alphanumeric() || ".-_ ".contains(c);
    let mut cleaned: String = base
        .chars()
        .map(|c| if allowed(c) { c } else { '_' })
        .collect();

    if cleaned.is_empty() {
        cleaned = "uploaded_file".to_string();
    }

    // Limit length (filesystem-friendly)
    if cleaned.len() > 255 {
        cleaned.truncate(255);
    }

    cleaned
}

fn resolve_unique_path(upload_dir: &str, name: &str) -> (String, PathBuf) {
    let cleaned = sanitize_filename(name);
    let stem = Path::new(&cleaned)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("file")
        .to_string();
    let ext = Path::new(&cleaned)
        .extension()
        .and_then(|s| s.to_str())
        .map(|s| s.to_string());

    let mut candidate = cleaned.clone();
    let mut path = PathBuf::from(upload_dir).join(&candidate);

    // If it already exists, append a timestamp
    while path.exists() {
        let ts = SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap()
            .as_millis();
        candidate = if let Some(ref e) = ext {
            format!("{}-{}.{}", stem, ts, e)
        } else {
            format!("{}-{}", stem, ts)
        };
        path = PathBuf::from(upload_dir).join(&candidate);
    }

    (candidate, path)
}

/// --- Handlers -----------------------------------------------------------

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

pub async fn upload_file(
    Query(params): Query<ListFilesQuery>,
    headers: HeaderMap,
    mut multipart: Multipart,
) -> Response {
    let upload_dir = params.dir.unwrap_or_else(|| UPLOAD_DIR.to_string());

    // Security check: Ensure upload_dir is within UPLOAD_DIR and doesn't contain ".."
    if !upload_dir.starts_with(UPLOAD_DIR) || upload_dir.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({
                "status": "error",
                "message": "❌ Access forbidden: You can only upload files inside the upload directory"
            })),
        )
            .into_response();
    }

    while let Ok(Some(mut field)) = multipart.next_field().await {
        let name = field.name().unwrap_or("").to_string();

        if name == "file" {
            // Prefer query.name -> X-Original-Filename header -> multipart filename
            let raw_name = params
                .name
                .clone()
                .or_else(|| {
                    headers
                        .get("x-original-filename")
                        .and_then(|h| h.to_str().ok().map(|s| s.to_string()))
                })
                .or_else(|| field.file_name().map(|s| s.to_string()))
                .unwrap_or_else(|| "uploaded_file".to_string());

            let (final_filename, file_path) = resolve_unique_path(&upload_dir, &raw_name);

            println!(
                "⬇️  Streaming start (Multipart): {} -> {}",
                raw_name, final_filename
            );

            if let Err(e) = tokio::fs::create_dir_all(&upload_dir).await {
                return (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(serde_json::json!({
                        "status": "error",
                        "message": format!("❌ Failed to create directory: {}", e)
                    })),
                )
                    .into_response();
            }

            let mut file = match File::create(&file_path).await {
                Ok(f) => f,
                Err(e) => {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Json(serde_json::json!({
                            "status": "error",
                            "message": format!("❌ Server Error: {}", e)
                        })),
                    )
                        .into_response()
                }
            };

            while let Ok(Some(chunk)) = field.chunk().await {
                if let Err(e) = file.write_all(&chunk).await {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Json(serde_json::json!({
                            "status": "error",
                            "message": format!("❌ Write Error: {}", e)
                        })),
                    )
                        .into_response();
                }
            }

            println!("✅ Upload complete: {}", final_filename);
            return Json(serde_json::json!({
                "status": "success",
                "message": format!("✅ Upload complete: {}", final_filename),
                "filename": final_filename
            }))
            .into_response();
        }
    }

    (
        StatusCode::BAD_REQUEST,
        Json(serde_json::json!({
            "status": "error",
            "message": "❌ No file field found in request"
        })),
    )
        .into_response()
}

pub async fn upload_raw(
    Query(params): Query<ListFilesQuery>,
    headers: HeaderMap,
    body: Body,
) -> Response {
    let upload_dir = params.dir.unwrap_or_else(|| UPLOAD_DIR.to_string());

    if !upload_dir.starts_with(UPLOAD_DIR) || upload_dir.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({
                "status": "error",
                "message": "❌ Access forbidden"
            })),
        )
            .into_response();
    }

    // Prefer query.name -> x-original-filename -> x-file-name -> fallback
    let raw_name = params
        .name
        .clone()
        .or_else(|| {
            headers
                .get("x-original-filename")
                .and_then(|h| h.to_str().ok().map(|s| s.to_string()))
        })
        .or_else(|| {
            headers
                .get("x-file-name")
                .and_then(|h| h.to_str().ok().map(|s| s.to_string()))
        })
        .unwrap_or_else(|| "uploaded_file".to_string());

    let (final_filename, file_path) = resolve_unique_path(&upload_dir, &raw_name);

    println!(
        "⬇️  Streaming start (Raw): {} -> {}",
        raw_name, final_filename
    );

    if let Err(e) = tokio::fs::create_dir_all(&upload_dir).await {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"status": "error", "message": e.to_string()})),
        )
            .into_response();
    }

    let file = match File::create(&file_path).await {
        Ok(f) => f,
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": e.to_string()})),
            )
                .into_response()
        }
    };

    let mut writer = BufWriter::new(file);
    let mut body_reader = StreamReader::new(
        body.into_data_stream()
            .map(|res| res.map_err(|err| std::io::Error::new(std::io::ErrorKind::Other, err))),
    );

    if let Err(e) = tokio::io::copy(&mut body_reader, &mut writer).await {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"status": "error", "message": e.to_string()})),
        )
            .into_response();
    }

    if let Err(e) = writer.flush().await {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"status": "error", "message": e.to_string()})),
        )
            .into_response();
    }

    println!("✅ Streaming complete: {}", final_filename);
    Json(serde_json::json!({
        "status": "success",
        "message": format!("✅ Upload complete: {}", final_filename),
        "filename": final_filename
    }))
    .into_response()
}

pub async fn stream_upload(
    Query(params): Query<ListFilesQuery>,
    headers: HeaderMap,
    mut multipart: Multipart,
) -> Response {
    let upload_dir = params.dir.unwrap_or_else(|| UPLOAD_DIR.to_string());

    if !upload_dir.starts_with(UPLOAD_DIR) || upload_dir.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({"status": "error", "message": "❌ Access forbidden"})),
        )
            .into_response();
    }

    while let Ok(Some(field)) = multipart.next_field().await {
        let name = field.name().unwrap_or("").to_string();

        if name == "file" {
            // Prefer params.name -> header -> multipart filename
            let raw_name = params
                .name
                .clone()
                .or_else(|| {
                    headers
                        .get("x-original-filename")
                        .and_then(|h| h.to_str().ok().map(|s| s.to_string()))
                })
                .or_else(|| field.file_name().map(|s| s.to_string()))
                .unwrap_or_else(|| "uploaded_file".to_string());

            let (final_filename, file_path) = resolve_unique_path(&upload_dir, &raw_name);

            println!(
                "⬇️  Advanced Streaming start: {} -> {}",
                raw_name, final_filename
            );

            if let Err(e) = tokio::fs::create_dir_all(&upload_dir).await {
                return (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(serde_json::json!({"status": "error", "message": e.to_string()})),
                )
                    .into_response();
            }

            let file = match File::create(&file_path).await {
                Ok(f) => f,
                Err(e) => {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Json(serde_json::json!({"status": "error", "message": e.to_string()})),
                    )
                        .into_response()
                }
            };

            let mut writer = BufWriter::new(file);
            let body_with_io_error = field
                .map(|res| res.map_err(|err| std::io::Error::new(std::io::ErrorKind::Other, err)));
            let mut body_reader = StreamReader::new(body_with_io_error);

            if let Err(e) = tokio::io::copy(&mut body_reader, &mut writer).await {
                return (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(serde_json::json!({"status": "error", "message": e.to_string()})),
                )
                    .into_response();
            }

            if let Err(e) = writer.flush().await {
                return (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(serde_json::json!({"status": "error", "message": e.to_string()})),
                )
                    .into_response();
            }

            println!("✅ Advanced Streaming complete: {}", final_filename);
            return Json(serde_json::json!({
                "status": "success",
                "message": format!("✅ Upload complete: {}", final_filename),
                "filename": final_filename
            }))
            .into_response();
        }
    }

    (
        StatusCode::BAD_REQUEST,
        Json(serde_json::json!({"status": "error", "message": "❌ No file field found"})),
    )
        .into_response()
}

pub async fn create_folder(Json(payload): Json<CreateFolderRequest>) -> Response {
    let full_path = PathBuf::from(UPLOAD_DIR).join(&payload.path);

    // Security check
    if !full_path.starts_with(UPLOAD_DIR) || payload.path.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({"status": "error", "message": "❌ Access forbidden"})),
        )
            .into_response();
    }

    match tokio::fs::create_dir_all(&full_path).await {
        Ok(_) => Json(serde_json::json!({
            "status": "success",
            "message": format!("✅ Folder created: {}", payload.path)
        }))
        .into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"status": "error", "message": e.to_string()})),
        )
            .into_response(),
    }
}

pub async fn create_file(Json(payload): Json<CreateFileRequest>) -> Response {
    let full_path = PathBuf::from(UPLOAD_DIR).join(&payload.path);

    // Security check
    if !full_path.starts_with(UPLOAD_DIR) || payload.path.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({"status": "error", "message": "❌ Access forbidden"})),
        )
            .into_response();
    }

    // Ensure parent directory exists
    if let Some(parent) = full_path.parent() {
        if let Err(e) = tokio::fs::create_dir_all(parent).await {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": e.to_string()})),
            )
                .into_response();
        }
    }

    let mut file = match File::create(&full_path).await {
        Ok(f) => f,
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": e.to_string()})),
            )
                .into_response()
        }
    };

    if let Some(content) = payload.content {
        if let Err(e) = file.write_all(content.as_bytes()).await {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": e.to_string()})),
            )
                .into_response();
        }
    }

    Json(serde_json::json!({
        "status": "success",
        "message": format!("✅ File created: {}", payload.path)
    }))
    .into_response()
}

pub async fn delete_item(Json(payload): Json<DeleteRequest>) -> Response {
    let full_path = PathBuf::from(UPLOAD_DIR).join(&payload.path);

    // Security check
    if !full_path.starts_with(UPLOAD_DIR) || payload.path.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({"status": "error", "message": "❌ Access forbidden"})),
        )
            .into_response();
    }

    if !full_path.exists() {
        return (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({"status": "error", "message": "❌ Path not found"})),
        )
            .into_response();
    }

    let result = if full_path.is_dir() {
        tokio::fs::remove_dir_all(&full_path).await
    } else {
        tokio::fs::remove_file(&full_path).await
    };

    match result {
        Ok(_) => Json(serde_json::json!({
            "status": "success",
            "message": format!("✅ Deleted: {}", payload.path)
        }))
        .into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"status": "error", "message": e.to_string()})),
        )
            .into_response(),
    }
}

pub async fn download_file(Query(params): Query<DownloadQuery>) -> Response {
    let full_path = PathBuf::from(UPLOAD_DIR).join(&params.path);

    // Security check
    if !full_path.starts_with(UPLOAD_DIR) || params.path.contains("..") {
        return (
            StatusCode::FORBIDDEN,
            Json(serde_json::json!({"status": "error", "message": "❌ Access forbidden"})),
        )
            .into_response();
    }

    if !full_path.exists() || full_path.is_dir() {
        return (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({"status": "error", "message": "❌ File not found"})),
        )
            .into_response();
    }

    let file = match File::open(&full_path).await {
        Ok(f) => f,
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": e.to_string()})),
            )
                .into_response()
        }
    };

    // Stream the file
    let stream = tokio_util::io::ReaderStream::new(file);
    let body = Body::from_stream(stream);

    let filename = full_path
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or("download");

    Response::builder()
        .header(
            "Content-Disposition",
            format!("attachment; filename=\"{}\"", filename),
        )
        .header("Content-Type", "application/octet-stream")
        .body(body)
        .unwrap()
        .into_response()
}
