# aaxion API — Quick Reference ⚡️

A concise, developer-friendly reference for the aaxion file-service API. Use the examples below to interact with a local server.

Base URL (local): `http://localhost:8080/`

---

## 🔐 Authentication

Most endpoints require authentication. You must obtain a token and include it in the `Authorization` header.

**Header format:**
`Authorization: Bearer <your_token>`

**Special Case:**

- `/files/thumbnail`: Supports passing the token via query parameter `?tkn=<token>` to allow loading images in `<img>` tags.

---

## Quick examples

- Login:

  ```bash
  curl -X POST -d '{"username":"your_user","password":"your_pass"}' "http://localhost:8080/auth/login"
  ```

- Logout:

  ```bash
  curl -H "Authorization: Bearer $TOKEN" -X POST "http://localhost:8080/auth/logout"
  ```

- View files:

  ```bash
  curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/files/view?dir=/home/swap/documents"
  ```

- Upload (single file):

  ```bash
  curl -H "Authorization: Bearer $TOKEN" -F "file=@/path/to/file" "http://localhost:8080/files/upload?dir=/home/swap/documents"
  ```

---

## Endpoints

### 👤 User Management

#### Register (Initial Setup)

Endpoint:

```http
POST /auth/register
```

- Description: Register the first user. Fails if a user already exists.
- Body: `{"username": "...", "password": "..."}`
- Response: HTTP 201 Created.

#### Login

Endpoint:

```http
POST /auth/login
```

- Description: Authenticate and receive a session token.
- Body: `{"username": "...", "password": "..."}`
- Response: `{"token": "..."}`

#### Logout

Endpoint:

```http
POST /auth/logout
```

- Description: Invalidate the current session token.
- Header: `Authorization: Bearer <token>`

---

### 📁 View Files and Folders

Endpoint:

```http
GET /files/view?dir={directory_path}
```

- **Requires Auth**: Yes
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

- **Requires Auth**: Yes
- Description: Create a new directory at the specified path.
- Parameters:
  - `path` (string, required): Target directory path (inside monitored root).
- Success response: empty body with HTTP `201 Created`.

Example:

```bash
curl -H "Authorization: Bearer $TOKEN" -X POST "http://localhost:8080/files/create-directory?path=/home/swap/new_folder"
```

---

### 📤 Upload File (single request)

Endpoint:

```http
POST /files/upload?dir={directory_path}
```

- **Requires Auth**: Yes
- Description: Upload a file via multipart form-data.
- Parameters:
  - `dir` (string, required): Destination directory.
  - Body: `multipart/form-data` with a `file` field.
- Success: HTTP `201 Created`.

Example:

```bash
curl -H "Authorization: Bearer $TOKEN" -F "file=@/tmp/example.txt" "http://localhost:8080/files/upload?dir=/home/swap/documents"
```

---

### ⚙️ Chunked Upload (for large files)

All chunked upload endpoints require `Authorization: Bearer <token>`.

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
curl -H "Authorization: Bearer $TOKEN" -X POST "http://localhost:8080/files/upload/chunk/start?filename=largeFile.zip"
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
curl -H "Authorization: Bearer $TOKEN" --data-binary @chunk0.bin "http://localhost:8080/files/upload/chunk?filename=largeFile.zip&chunk_index=0"
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
curl -H "Authorization: Bearer $TOKEN" -X POST "http://localhost:8080/files/upload/chunk/complete?filename=largeFile.zip&dir=/home/swap/documents"
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
- Response: JSON object containing the share link and token.
- Example response:

  ```json
  {
    "share_link": "/files/d/t/abcdefghijklmnopqrstuvwxyzABCDEF",
    "token": "abcdefghijklmnopqrstuvwxyzABCDEF"
  }
  ```

Use the token URL to download:

Endpoint:

```http
GET /files/d/t/{token}
```

- **Requires Auth**: No (Token acts as auth)
- Description: Download the file referenced by the one-time token.
- Note: Tokens are valid for one use only.

Example:

```bash
curl -O "http://localhost:8080/files/d/t/abcdefghijklmnopqrstuvwxyzABCDEF"
```

---

### 🖼️ Images & Thumbnails

#### View Full Image

Endpoint:

```http
GET /files/view-image?path={file_path}
```

- **Requires Auth**: Yes
- Description: Serve the raw image file directly.
- Features:
  - Supports client-side caching (7 days).
  - Handles correct content-type automatically.
- Query: `path` (string, required).

Example:

```bash
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/files/view-image?path=/home/swap/photos/vacation.jpg"
```

#### Get Thumbnail

Endpoint:

```http
GET /files/thumbnail?path={file_path}
```

- **Requires Auth**: Yes (Header or Query Param)
- Description: Get a resized (max 200px) JPEG thumbnail of an image.
- Features: Server-side caching of generated thumbnails.
- Query:
  - `path` (string, required).
  - `tkn` (string, optional): Auth token, for use in `<img>` tags.

Example (Header):

```bash
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/files/thumbnail?path=/home/swap/photos/vacation.jpg"
```

Example (Query Param):

```bash
curl "http://localhost:8080/files/thumbnail?path=/home/swap/photos/vacation.jpg&tkn=$TOKEN"
```

---

### System Info

Endpoint:

```http
GET /api/system/get-root-path
```

- **Requires Auth**: Yes
- Description: Retrieve the monitored root directory path.
- Response: JSON object with `root_path` field.
  Example response:

```json
{
  "root_path": "/home/swap"
}
```

## Notes & best practices

- All paths must be under the monitored root (e.g., `/home/swap/*`).
- Ensure filesystem permissions allow the server process to read/write the target locations.
- For large uploads prefer the chunked flow; keep chunks <= 90MB.

---
