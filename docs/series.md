# Series API Reference 🎬

Documentation for the Series and Episodes management API endpoints.

Base URL (local): `http://localhost:8080/`

---

## 🔐 Authentication

All endpoints require authentication (except where noted otherwise for testing/token scenarios).
Include the `Authorization` header:

```
Authorization: Bearer <your_token>
```

---

## 📺 Series Management

### List All Series

Endpoint:

```http
GET /api/series/list
```

- **Requires Auth**: Yes
- **Description**: Returns a list of all available series.
- **Response**: JSON array of series objects.

Example Response:

```json
[
  {
    "id": 1,
    "title": "Breaking Bad",
    "description": "A high school chemistry teacher turned methamphetamine producer.",
    "created_at": "2023-10-27T10:00:00Z"
  }
]
```

### Search Series

Endpoint:

```http
GET /api/series/search?q={query}
```

- **Requires Auth**: Yes
- **Description**: Search for series by title or description.
- **Parameters**:
  - `q` (string, required): The search query.
- **Response**: JSON array of matching series.

### Add Series

Endpoint:

```http
POST /api/series/add
```

- **Requires Auth**: Yes
- **Description**: Create a new series entry.
- **Body**: JSON object.

Example Request:

```json
{
  "title": "Stranger Things",
  "description": "When a young boy disappears, his mother, a police chief and his friends must confront terrifying supernatural forces."
}
```

- **Response**: HTTP 201 Created.

### Edit Series

Endpoint:

```http
PUT /api/series/edit
```

- **Requires Auth**: Yes
- **Description**: Update an existing series information.
- **Body**: JSON object.

Example Request:

```json
{
  "id": 1,
  "title": "Stranger Things (Season 1)",
  "description": "Updated description..."
}
```

- **Response**: HTTP 200 OK.

---

## 🎞️ Episode Management

### List Episodes

Endpoint:

```http
GET /api/series/episodes/list?series_id={id}
```

- **Requires Auth**: Yes
- **Description**: List all episodes for a specific series, ordered by season and episode number.
- **Parameters**:
  - `series_id` (int, required): ID of the series.
- **Response**: JSON array of episode objects.

Example Response:

```json
[
  {
    "id": 101,
    "series_id": 1,
    "season_number": 1,
    "episode_number": 1,
    "title": "The Vanishing of Will Byers",
    "description": "On his way home from a friend's house, young Will sees something terrifying...",
    "file_path": "/path/to/S01E01.mp4",
    "size": 104857600,
    "mime_type": "video/mp4",
    "created_at": "2023-10-27T10:05:00Z"
  }
]
```

### Add Episode

Endpoint:

```http
POST /api/series/episodes/add
```

- **Requires Auth**: Yes
- **Description**: Add an episode to a series.
- **Body**: JSON object.

Example Request:

```json
{
  "series_id": 1,
  "file_id": 123,
  "file_path": "/media/series/stranger_things/s01e01.mp4",
  "season_number": 1,
  "episode_number": 1,
  "title": "Chapter One",
  "description": "The pilot episode."
}
```

- **Response**: HTTP 201 Created.

### Stream Episode

Endpoint:

```http
GET /api/stream/episode?id={episode_id}
```

- **Requires Auth**: Yes
- **Description**: Stream the episode video content. Supports HTTP Range requests for seeking.
- **Parameters**:
  - `id` (int, required): The ID of the episode to stream.
- **Response**: Video stream (binary content).

---
