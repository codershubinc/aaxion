use axum::{extract::DefaultBodyLimit, http::Method, middleware, Router};
use tower_http::cors::{Any, CorsLayer};

use crate::{middlewares::auth_middleware::require_auth, routes::file_routes};

pub fn create_router() -> Router {
    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods([
            Method::GET,
            Method::POST,
            Method::PUT,
            Method::DELETE,
            Method::OPTIONS,
        ])
        .allow_headers(Any);

    Router::new()
        .nest(
            "/api/files",
            file_routes::routes()
                .layer(middleware::from_fn(require_auth))
                .layer(DefaultBodyLimit::disable()),
        )
        .layer(cors)
}
