use axum::{extract::DefaultBodyLimit, middleware, Router};

use crate::{middlewares::auth_middleware::require_auth, routes::file_routes};

pub fn create_router() -> Router {
    Router::new().nest(
        "/api/files",
        file_routes::routes()
            .layer(middleware::from_fn(require_auth))
            .layer(DefaultBodyLimit::disable()),
    )
}
