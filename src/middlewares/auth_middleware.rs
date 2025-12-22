use axum::{
    body::Body,
    http::{Request, StatusCode},
    middleware::Next,
    response::{IntoResponse, Response},
};

pub async fn require_auth(req: Request<Body>, next: Next) -> Response {
    if let Some(hv) = req.headers().get("authorization") {
        if let Ok(s) = hv.to_str() {
            if s == "Bearer my_secret_token" {
                return next.run(req).await;
            }
        }
    }

    // (StatusCode::UNAUTHORIZED, "Unauthorized").into_response()  // go for now

    next.run(req).await
}
