# `upload_file` — Line-by-line explanation ✅

**Purpose:** Handles multipart `POST` uploads and writes the incoming file stream to disk under a configured upload directory. This documents `pub async fn upload_file(mut multipart: Multipart) -> impl IntoResponse` from `src/handlers/file_handler.rs` and explains each line.

---

```rust
pub async fn upload_file(mut multipart: Multipart) -> impl IntoResponse {
    // We loop through ALL fields until we find the one we want
    while let Ok(Some(mut field)) = multipart.next_field().await {
        let name = field.name().unwrap_or("").to_string();

        if name == "file" {
            // FIX: Handle missing filename safely
            let file_name = match field.file_name() {
                Some(n) => n.to_string(),
                None => continue, // Skip if it's not a file
            };

            // SECURITY FIX: Sanitize filename (prevent ../../ attacks)
            let safe_filename = Path::new(&file_name).file_name().unwrap().to_string_lossy();
            let file_path = PathBuf::from(UPLOAD_DIR).join(safe_filename.as_ref());

            println!("⬇️  Streaming start: {}", safe_filename);

            let mut file = match File::create(&file_path).await {
                Ok(f) => f,
                Err(e) => {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("❌ Server Error: {}", e),
                    )
                        .into_response()
                }
            };

            while let Ok(Some(chunk)) = field.chunk().await {
                if let Err(e) = file.write_all(&chunk).await {
                    return (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("❌ Write Error: {}", e),
                    )
                        .into_response();
                }
            }

            println!("✅ Streaming complete: {}", safe_filename);
            return Html(format!("✅ Upload complete: {}", safe_filename)).into_response();
        }
    }

    // FIX: Provide a return value if the loop finishes without finding a file
    (StatusCode::BAD_REQUEST, "❌ No file field found in request").into_response()
}
```

---

Line-by-line explanation:

- `pub async fn upload_file(mut multipart: Multipart) -> impl IntoResponse {`

  - `pub`: publicly accessible handler used by the router.
  - `async`: uses async I/O to stream and write file chunks.
  - `mut multipart: Multipart`: Axum's multipart extractor for `multipart/form-data` requests; `mut` because we iterate and consume fields.
  - `-> impl IntoResponse`: returns a type that can be converted into an Axum `Response` (flexible: HTML, JSON, or error tuple).

- `while let Ok(Some(mut field)) = multipart.next_field().await {`

  - Iterates over all fields in the multipart body; `next_field().await` returns each field or `Ok(None)` when exhausted. Using `while let Ok(Some(...))` ignores non-fatal errors and stops on errors.

- `let name = field.name().unwrap_or("").to_string();`

  - Fetches the form field name (the `name` attribute in the HTML form). If not present, defaults to an empty string.

- `if name == "file" {` ... `}`

  - The handler expects the file field to be named `file` (common convention). If a different name is used, the handler will skip until it finds `file`.

- `let file_name = match field.file_name() { Some(n) => n.to_string(), None => continue, };`

  - `field.file_name()` returns the filename supplied by the client (may be absent for non-file fields).
  - If filename is missing, the code `continue`s to skip to the next field (safe fallback).

- `let safe_filename = Path::new(&file_name).file_name().unwrap().to_string_lossy();`

  - **Sanitization step:** `Path::new(...).file_name()` extracts only the _final_ component of the filename (for example, it strips any path like `../../etc/passwd` to just `passwd`) which prevents simple path traversal attacks.
  - `to_string_lossy()` returns a `Cow<str>` converting non-UTF8 bytes into a presentation-friendly string.

- `let file_path = PathBuf::from(UPLOAD_DIR).join(safe_filename.as_ref());`

  - Joins the sanitized filename to the server's configured `UPLOAD_DIR` to obtain the final destination path.
  - This ensures uploads always go into the configured upload directory, not a client-controlled location.

- `println!("⬇️  Streaming start: {}", safe_filename);`

  - Logs that a streaming upload has started for monitoring/debugging.

- `let mut file = match File::create(&file_path).await { Ok(f) => f, Err(e) => { return (StatusCode::INTERNAL_SERVER_ERROR, format!("❌ Server Error: {}", e),).into_response() } };`

  - Attempts to create a file asynchronously; on success obtains an async `File` handle.
  - On failure (permissions, disk full, invalid path), returns a 500 Internal Server Error with the formatted underlying error message. This early return stops processing further fields.

- `while let Ok(Some(chunk)) = field.chunk().await { if let Err(e) = file.write_all(&chunk).await { return (StatusCode::INTERNAL_SERVER_ERROR, format!("❌ Write Error: {}", e),).into_response(); } }`

  - Reads the field body in chunks (`field.chunk().await`) to support streaming large files without holding everything in memory.
  - Each chunk is written asynchronously to the created file with `file.write_all(&chunk).await`.
  - On write failure, returns 500 with write error details.

- `println!("✅ Streaming complete: {}", safe_filename);`

  - Logs that streaming finished successfully.

- `return Html(format!("✅ Upload complete: {}", safe_filename)).into_response();`

  - Returns an HTML response confirming the upload. `Html(...)` wraps a string into an `axum::response::Html` type which sets the Content-Type header.

- After the loop: `(StatusCode::BAD_REQUEST, "❌ No file field found in request").into_response()`
  - If the handler finished iterating fields without finding a `file` field, it returns HTTP 400 Bad Request with a short text message.

---

Notes & recommendations:

- The sanitization via `file_name()` is effective for typical attacks but consider further validation:

  - Normalize and canonicalize the final `file_path` and ensure it starts with `UPLOAD_DIR` to be extra safe.
  - Enforce filename rules (e.g., allowed characters, length limits) to prevent weird filenames on the filesystem.

- Concurrency & atomicity:

  - Two concurrent uploads with the same filename will clobber each other. Consider creating unique filenames (e.g., prefix with a UUID or timestamp) or using locking.

- Error transparency:

  - Right now server returns internal error messages to clients (with `format!("{}", e)`) — be mindful of leaking server internals in production; prefer logging the detailed error and returning a user-friendly message.

- Streaming and memory:

  - The approach avoids loading the whole file into memory and is appropriate for large files.

- Tests: Add tests for missing `file_name`, path traversal attempts in the filename, large file streaming, and write failure handling.

---

If you'd like, I can:

- Add a sample upload HTML form and a cURL example to the docs, or
- Implement automatic filename deduplication (append a counter/UUID) and add tests.
