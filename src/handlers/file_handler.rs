use axum::{extract::Query, response::IntoResponse, Json};
use serde::Deserialize;

use crate::models::FileInfo;
#[derive(Deserialize)]
pub struct ListFilesQuery {
    pub dir: Option<String>,
}

pub async fn list_files(Query(params): Query<ListFilesQuery>) -> impl IntoResponse {
    let dir = params
        .dir
        .unwrap_or_else(|| "/home/swap/aaxion/".to_string());
    let mut entries = match tokio::fs::read_dir(&dir).await {
        Ok(e) => e,
        Err(_) => return Json(vec![]),
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
                    path: dir.clone(),
                    raw_path: path.to_string_lossy().to_string(),
                });
            }
        }
    }

    Json(files)
}
