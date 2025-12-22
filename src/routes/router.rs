use axum::{middleware, Router};

use crate::{middlewares::auth_middleware::require_auth, routes::file_routes};

pub fn create_router() -> Router {
    Router::new().nest(
        "/api/files",
        file_routes::routes().layer(middleware::from_fn(require_auth)),
    )
}
