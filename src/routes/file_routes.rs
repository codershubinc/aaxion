use axum::{
    routing::{get, post},
    Router,
};

use crate::handlers;

pub fn routes() -> Router {
    Router::new()
        .route("/list", get(handlers::list_files))
        .route("/upload", post(handlers::upload_file))
        .route("/upload-raw", axum::routing::put(handlers::upload_raw))
        .route("/stream-upload", post(handlers::stream_upload))
        .route("/create-folder", post(handlers::create_folder))
        .route("/create-file", post(handlers::create_file))
        .route("/delete", post(handlers::delete_item))
        .route("/download", get(handlers::download_file))
}
