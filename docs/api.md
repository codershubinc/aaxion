# aaxion API — Quick Reference ⚡️

A concise, developer-friendly reference for the aaxion file-service API. Use the examples below to interact with a local server.

Base URL (local): `http://localhost:8080/`

---

## Quick examples

- View files (curl):

  ```bash
  curl "http://localhost:8080/files/view?dir=/home/swap/documents"
  ```

- Upload (single file):

  ```bash
  curl -F "file=@/path/to/file" "http://localhost:8080/files/upload?dir=/home/swap/documents"
  ```

---

## Endpoints

### 📁 View Files and Folders

Endpoint:

```http
GET /files/view?dir={directory_path}
```

- Description: Return the contents of a directory.
- Parameters:
  - `dir` (string, required): Path of the directory to list (must be inside the monitored root).
- Response: JSON array of file/folder objects.
- Example response:

  ```json
  [
    {
      "name": "Quazaar",
      "is_dir": true,
      "size": 4096,
      "path": "/home/swap/Github",
      "raw_path": "/home/swap/Github/Quazaar"
    }
  ]
  ```

⚠️ Warning: `dir` must start within the monitored root (e.g., `/home/swap/`). Requests outside the root will be rejected with a "Suspicious path detected" error.

---

### ✨ Create Directory

Endpoint:

```http
POST /files/create-directory?path={directory_path}
```

- Description: Create a new directory at the specified path.
- Parameters:
  - `path` (string, required): Target directory path (inside monitored root).
- Success response: empty body with HTTP `201 Created`.

Example:

```bash
curl -X POST "http://localhost:8080/files/create-directory?path=/home/swap/new_folder"
```

---

### 📤 Upload File (single request)

Endpoint:

```http
POST /files/upload?dir={directory_path}
```

- Description: Upload a file via multipart form-data.
- Parameters:
  - `dir` (string, required): Destination directory.
  - Body: `multipart/form-data` with a `file` field.
- Success: HTTP `201 Created`.

Example:

```bash
curl -F "file=@/tmp/example.txt" "http://localhost:8080/files/upload?dir=/home/swap/documents"
```

---

### ⚙️ Chunked Upload (for large files)

1. Start session

Endpoint:

```http
POST /files/upload/chunk/start?filename={filename}
```

- Description: Initialize a chunked upload session.
- Query: `filename` (required).
- Response: `Upload initialized` (or JSON with session info).

Example:

```bash
curl -X POST "http://localhost:8080/files/upload/chunk/start?filename=largeFile.zip"
```

2. Upload chunk

Endpoint:

```http
POST /files/upload/chunk?filename={filename}&chunk_index={index}
```

- Description: Upload a single chunk. Body is raw binary (NOT multipart/form-data).
- Query:
  - `filename` (required)
  - `chunk_index` (int, required) — start from `0` and increment by 1.
- Note: Upload chunks in order. Keep chunks <= 90MB for reliability.

Example:

```bash
curl --data-binary @chunk0.bin "http://localhost:8080/files/upload/chunk?filename=largeFile.zip&chunk_index=0"
```

3. Complete upload

Endpoint:

```http
POST /files/upload/chunk/complete?filename={filename}&dir={directory_path}
```

- Description: Merge uploaded chunks into the final file and save to `dir`.
- Query: `filename`, `dir` (required).
- Response: `File merged successfully`.

Example:

```bash
curl -X POST "http://localhost:8080/files/upload/chunk/complete?filename=largeFile.zip&dir=/home/swap/documents"
```

---

### 🔗 Temporary File Sharing

Request a temporary link (server returns a token / short URL):

Endpoint:

```http
GET /files/d/r?file_path={file_path}
```

- Description: Generate a one-time temporary link for a file.
- Query: `file_path` (required).
- Example response: `/files/d/t/abcdefghijklmnopqrstuvwxyzABCDEF`

Use the token URL to download:

Endpoint:

```http
GET /files/d/t/{token}
```

- Description: Download the file referenced by the one-time token.
- Note: Tokens are valid for one use only.

Example:

```bash
curl -O "http://localhost:8080/files/d/t/abcdefghijklmnopqrstuvwxyzABCDEF"
```

---

## Notes & best practices

- All paths must be under the monitored root (e.g., `/home/swap/*`).
- Ensure filesystem permissions allow the server process to read/write the target locations.
- For large uploads prefer the chunked flow; keep chunks <= 90MB.

--- 
