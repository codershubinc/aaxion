use axum::Router;

use crate::routes::file_routes;

pub fn create_router() -> Router {
    Router::new().nest("/api/files", file_routes::routes())
}
