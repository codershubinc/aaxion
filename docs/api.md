# Aaxion API — Deep Reference ⚡️

A complete, developer-friendly reference for the Aaxion REST API (v1). All endpoints are prefixed with `/api/v1/`.

Base URL (local): `http://localhost:8080/api/v1/`

---

## 🔐 Authentication

Most endpoints require authentication via a Bearer token. 

**Header format:**
`Authorization: Bearer <your_token>`

**Special Case:**
Endpoints like `/images/thumbnail`, `/files/download`, and streaming endpoints support passing the token via query parameter `?tkn=<token>` to allow loading in `<img>`, `<video>`, or `<a>` tags.

---

## 👤 User & Auth Management

### Register Admin
**`POST /api/v1/auth/register`** (Public)
- **Description:** Register the initial admin user. Fails if a user already exists.
- **Request Body:** `{"username": "admin", "password": "mypassword"}`
- **Response (201):** `{"message": "User created"}`

### Login
**`POST /api/v1/auth/login`** (Public)
- **Description:** Authenticate and receive session tokens.
- **Request Body:** `{"username": "admin", "password": "mypassword"}`
- **Response (200):**
  ```json
  {
    "token": "<session_token>",
    "access_token": "<access_token>",
    "device_info": { "id": "...", "name": "..." }
  }
  ```

### Generate Access Token
**`POST /api/v1/auth/token/generate`** (Auth Required)
- **Description:** Generate a short-lived access token. Requires primary session token.
- **Response (201):** `{"token": "<new_token>"}`

---

## 📁 Files & Directories

### View Directory
**`GET /api/v1/files/view?dir={path}`** (Auth Required)
- **Description:** List contents of a directory.
- **Response (200):**
  ```json
  [
    {
      "name": "Documents",
      "is_dir": true,
      "size": 4096,
      "path": "/home/user",
      "raw_path": "/home/user/Documents"
    }
  ]
  ```

### Unzip Archive
**`POST /api/v1/files/unzip`** (Auth Required)
- **Description:** Extract a ZIP archive.
- **Request Body:** `{"zip_path": "/path/archive.zip", "dest_dir": "/path/out"}`
- **Response (200):** `{"message": "Archive extracted successfully", "dest": "/path/out"}`

---

## ⚙️ Chunked Uploads (Large Files)

Requires Auth. For files larger than 100MB.

1. **`POST /api/v1/files/upload/chunk/start?filename=big.zip`**
   - **Response:** `Upload initialized`
2. **`POST /api/v1/files/upload/chunk?filename=big.zip&chunk_index=0`**
   - **Body:** Raw binary chunk. Keep under 90MB.
   - **Response:** `Chunk received`
3. **`POST /api/v1/files/upload/chunk/complete?filename=big.zip&dir=/dest`**
   - **Response:** `File merged successfully`

---

## 🔗 Temporary Sharing & Anonymous

### Generate Temp Link
**`GET /api/v1/share/temp/request?file_path={path}`** (Auth Required)
- **Response (200):**
  ```json
  {
    "share_link": "/api/v1/share/temp/abc...",
    "token": "abc..."
  }
  ```

### Anonymous Upload Token Generation
**`POST /api/v1/anonymous/token/generate?target_dir=/up&max_uploads=5`** (Auth Required)
- **Response (200):**
  ```json
  {
    "token": "abc123",
    "upload_url": "<host>/upload?token=abc123",
    "target_dir": "/up",
    "max_uploads": 5,
    "expiry_hours": 24,
    "max_file_size": 11811160064
  }
  ```

### Anonymous File Upload
**`POST /api/v1/anonymous/upload?token=abc123`** (Public)
- **Body:** `multipart/form-data` with `file`
- **Response (201):** `{"message": "Upload successful", "uploads_remaining": 4}`

---

## 🎵 Music

### Add Music
**`POST /api/v1/music/add?uri={youtube_url}`** (Public)
- **Description:** Download music via yt-dlp. URI passed via form value or query.
- **Response (202):** `{"status": "success", "message": "Queued 1 tracks", "count": 1}`

### Stream Music
**`GET /api/v1/music/stream?id=1`** (Public)
- **Response:** `<binary audio stream>`

### Music Stats
**`POST /api/v1/music/stats/play?track_id=1`** (Public)
**`POST /api/v1/music/stats/favorite?track_id=1&is_favorite=true`** (Public)

---

## 💻 System

### System Storage
**`GET /api/v1/system/storage`** (Auth Required)
- **Response (200):**
  ```json
  {
    "total": 512000000000,
    "used": 256000000000,
    "available": 256000000000,
    "usage_percentage": 50.0,
    "external_devices": [
      {
        "device": "/dev/sdb1",
        "mount_point": "/media/drive",
        "filesystem_type": "ext4",
        "total": 64000000000,
        "used": 32000000000,
        "available": 32000000000,
        "usage_percentage": 50.0
      }
    ]
  }
  ```
