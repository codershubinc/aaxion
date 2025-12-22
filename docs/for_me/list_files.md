# `list_files` — Line-by-line explanation ✅

**Purpose:** Lists files inside the configured upload directory and returns a JSON array of `FileInfo` objects (skips hidden files and prevents directory traversal). This documents the function `pub async fn list_files(Query(params): Query<ListFilesQuery>) -> Response` from `src/handlers/file_handler.rs` and explains each line and decision.

---

```rust
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
        }
    }

    Json(files).into_response()
}
```

---

Line-by-line explanation:

- `pub async fn list_files(Query(params): Query<ListFilesQuery>) -> Response {`

  - `pub`: function is public so other modules/routers can call it.
  - `async`: function is asynchronous; uses `await` for I/O (non-blocking).
  - `list_files(...)`: handler name used in routing.
  - `Query(params): Query<ListFilesQuery>`: uses Axum's `Query` extractor to deserialize URL query parameters into the `ListFilesQuery` struct (which has an optional `dir: Option<String>`). The pattern `Query(params)` moves the inner `params` out of the wrapper for convenience.
  - `-> Response`: returns an Axum `Response` type (allows different response kinds: JSON, status code, HTML).

- `let requested_dir = params.dir.unwrap_or_else(|| UPLOAD_DIR.to_string());`

  - Reads the optional `dir` query param; if absent, defaults to the constant `UPLOAD_DIR` (configured root for uploads).
  - `unwrap_or_else` is lazy, so `UPLOAD_DIR.to_string()` runs only if needed.

- `println!("📂 Listing files in directory: {}", requested_dir);`

  - Debug/log line to stdout (visible when running server). Helpful for tracking which directory is requested.

- `let path = Path::new(&requested_dir);`

  - Creates a `Path` reference for filesystem operations.

- `if !requested_dir.starts_with(UPLOAD_DIR) || requested_dir.contains("..") {` ... `}`

  - **Security check**: prevents traversal outside `UPLOAD_DIR` and blocks any `..` in the path.
  - `starts_with(UPLOAD_DIR)` ensures the requested path is under the allowed upload root.
  - `requested_dir.contains("..")` is a simple extra guard to block `..` occurrences; this prevents basic directory-traversal attempts.
  - If the check fails, returns an HTTP 403 Forbidden and a plain text message. `.into_response()` converts the tuple `(StatusCode, &str)` into a proper Axum `Response`.

- `let mut entries = match tokio::fs::read_dir(&path).await { Ok(e) => e, Err(_) => return Json(vec![] as Vec<FileInfo>).into_response(), };`

  - Attempts to read the directory asynchronously using Tokio's `read_dir`.
  - On success, returns a `ReadDir` stream in `entries` for iterating entries.
  - On error (e.g., directory doesn't exist or permission denied), the function returns an empty JSON array (`Vec<FileInfo>`) as the response. This is a graceful fallback so the endpoint still returns JSON in failure cases.

- `let mut files: Vec<FileInfo> = Vec::new();`

  - Prepares an empty vector to collect `FileInfo` objects (this type is defined in `crate::models::FileInfo`).

- `while let Ok(Some(entry)) = entries.next_entry().await {` ... `}`

  - Asynchronously iterates directory entries; `next_entry().await` gives `Ok(Some(entry))` when there is a next entry; `Ok(None)` means the dir is exhausted.

  - Inside the loop:

    - `let path = entry.path();` — the full path for the entry (file or directory).
    - `let size = entry.metadata().await.map(|m| m.len()).unwrap_or(0);` — fetches file metadata to get the size in bytes; on failure uses `0`.

    - `if let Some(name) = path.file_name().and_then(|n| n.to_str()) {` — extract a UTF-8 filename string; skip non-UTF8 names.

      - `if !name.starts_with(".") {` — skip hidden files (names that start with `.`)

        - `files.push(FileInfo { name: name.to_string(), is_dir: path.is_dir(), size, path: requested_dir.clone(), raw_path: path.to_string_lossy().to_string(), });` — create and push a `FileInfo`:
          - `name`: filename as `String`.
          - `is_dir`: whether the entry is a directory.
          - `size`: file size in bytes (0 for unreadable metadata).
          - `path`: the requested directory string (where the file was listed from).
          - `raw_path`: the platform-native path serialized to a String (lossy conversion for non-utf8 parts).

- `Json(files).into_response()`
  - Wraps the collected `files` in Axum's `Json` response type (serializes to JSON) and converts it into a `Response` to return to the client.

---

Notes & recommendations:

- The function uses simple checks to prevent traversal (starts_with + contains("..")), which is good but not perfect for every edge case. Consider canonicalizing the path with `std::fs::canonicalize` (or equivalent async approach) and verifying the canonicalized path is under `UPLOAD_DIR`.

- Returning an empty array on read errors hides the exact issue from the client (could be permission error vs non-existent directory). This is acceptable for a public API but consider returning a more specific error for debugging or admin endpoints.

- Performance: fetching metadata for each entry is required to get file size, but it adds an async await per file. For large directories you may want streaming responses or pagination.

- Tests: Add unit/integration tests for requests with `dir` query param, requests with traversal attempts, and empty directories.

---

If you want, I can also add sample cURL/JS examples and a short test harness for this handler.
