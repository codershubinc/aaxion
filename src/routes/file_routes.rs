use axum::{
    routing::{get, post},
    Router,
};

use crate::handlers;

pub fn routes() -> Router {
    Router::new()
        .route("/list", get(handlers::list_files))
        .route("/upload", post(handlers::upload_file))
}
