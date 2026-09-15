# Aaxion Series & Movies API Reference 🎬

*(Note: These endpoints are legacy and not actively developed, but are still present in v0.0.1-beta)*

Base URL: `http://localhost:8080/api/v1/`

---

## 🎥 Movies

### List Movies
**`GET /api/v1/movies`** (Auth Required)
- **Response:**
  ```json
  [
    {
      "id": 1,
      "title": "Interstellar",
      "file_path": "/path/to/movie.mp4",
      "file_id": 10,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
  ```

### Add Movie
**`POST /api/v1/movies/add`** (Auth Required)
- **Request Body:**
  ```json
  {
    "title": "Interstellar",
    "file_path": "/path/to/movie.mp4",
    "description": "Optional description",
    "poster_path": "/path/to/poster.jpg"
  }
  ```
- **Response (201):** Empty body.

### Edit Movie
**`PUT /api/v1/movies/edit`** (Auth Required)
- **Request Body:** `{"id": 1, "title": "Updated Title"}`
- **Response (200):** Empty body.

### Stream Movie
**`GET /api/v1/movies/stream?id=1`** (Auth Required)
- **Response:** Binary video stream (HTTP 206 Partial Content supported).

---

## 📺 Series & Episodes

### List Series
**`GET /api/v1/series`** (Auth Required)
- **Response:**
  ```json
  [
    {
      "id": 1,
      "title": "Breaking Bad",
      "description": "Series description",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
  ```

### Add Series
**`POST /api/v1/series/add`** (Auth Required)
- **Request Body:** `{"title": "Breaking Bad", "description": "Crime drama"}`
- **Response (201):** Empty body.

### Edit Series
**`PUT /api/v1/series/edit`** (Auth Required)
- **Request Body:** `{"id": 1, "title": "Updated Title"}`
- **Response (200):** Empty body.

### List Episodes
**`GET /api/v1/series/episodes?series_id=1`** (Auth Required)
- **Response:**
  ```json
  [
    {
      "id": 101,
      "series_id": 1,
      "season_number": 1,
      "episode_number": 1,
      "title": "Pilot",
      "file_path": "/path/to/episode.mp4",
      "size": 524288000,
      "mime_type": "video/mp4"
    }
  ]
  ```

### Add Episode
**`POST /api/v1/series/episodes/add`** (Auth Required)
- **Request Body:**
  ```json
  {
    "series_id": 1,
    "file_id": 10,
    "file_path": "/path/to/episode.mp4",
    "season_number": 1,
    "episode_number": 1,
    "title": "Pilot",
    "description": "Episode 1 description"
  }
  ```
- **Response (201):** Empty body.

### Stream Episode
**`GET /api/v1/series/episodes/stream?id=101`** (Auth Required)
- **Response:** Binary video stream.
