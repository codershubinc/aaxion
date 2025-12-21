use axum::{routing::get, Router};

use crate::handlers;

pub fn routes() -> Router {
    Router::new().route("/list", get(handlers::list_files))
}
