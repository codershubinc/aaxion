import Link from 'next/link';
import { Separator } from '@/components/ui/separator';
import { Button } from '@/components/ui/button';
import { CopyButton } from '@/components/copy-button';
import { ExpandableCode } from '@/components/expandable-code';

export default function Docs() {
  return (
    <div className="min-h-screen bg-[#050508] text-white selection:bg-indigo-500/30 relative overflow-hidden">
      <div className="absolute top-[-15%] left-[20%] w-[50%] h-[40%] rounded-full bg-indigo-700/10 blur-[150px] pointer-events-none" />
      <div className="absolute bottom-[10%] right-[-10%] w-[40%] h-[30%] rounded-full bg-cyan-700/8 blur-[130px] pointer-events-none" />

      <nav className="border-b border-white/5 sticky top-0 z-50 bg-[#050508]/70 backdrop-blur-2xl">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-14 items-center">
            <Link href="/" className="text-xl font-bold bg-gradient-to-r from-indigo-400 to-cyan-400 bg-clip-text text-transparent">Aaxion</Link>
            <div className="flex items-center gap-1">
              <div className="hidden lg:flex">
                <NavLink href="#authentication">Auth</NavLink>
                <NavLink href="#files">Files</NavLink>
                <NavLink href="#upload-download">Upload</NavLink>
                <NavLink href="#chunked">Chunks</NavLink>
                <NavLink href="#media">Media</NavLink>
                <NavLink href="#sharing">Sharing</NavLink>
                <NavLink href="#anonymous">Anonymous</NavLink>
                <NavLink href="#system">System</NavLink>
                <NavLink href="#movies">Movies</NavLink>
                <NavLink href="#series">Series</NavLink>
                <NavLink href="#music">Music</NavLink>
                <NavLink href="#websockets">WS</NavLink>
              </div>
              <Link href="/" className="ml-4">
                <Button variant="outline" size="sm" className="bg-white/5 border-white/10 hover:bg-white/10 hover:text-white text-xs h-8">&larr; Home</Button>
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-4xl mx-auto px-4 py-16 relative z-10">
        <h1 className="text-4xl md:text-5xl font-extrabold mb-3 tracking-tight">API Reference <span className="text-gray-600 font-normal text-2xl ml-2">v1</span></h1>
        <p className="text-base text-gray-500 mb-8 max-w-xl">Complete reference for the Aaxion REST API.</p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mb-16">
          <div className="bg-white/[0.025] border border-white/[0.06] rounded-xl p-4 flex items-center justify-between gap-4">
            <div>
              <p className="text-[10px] text-gray-600 uppercase tracking-widest font-semibold mb-1">Base URL</p>
              <code className="text-indigo-300 font-mono text-sm">http://localhost:8080/api/v1</code>
            </div>
            <CopyButton text="http://localhost:8080/api/v1" />
          </div>
          <div className="bg-white/[0.025] border border-white/[0.06] rounded-xl p-4 flex items-center justify-between gap-4">
            <div>
              <p className="text-[10px] text-gray-600 uppercase tracking-widest font-semibold mb-1">Auth Header</p>
              <code className="text-indigo-300 font-mono text-sm">Bearer &lt;token&gt;</code>
            </div>
            <CopyButton text="Authorization: Bearer <your_token>" />
          </div>
        </div>

        {/* ─── Authentication ─── */}
        <Section id="authentication" title="Authentication">
          <EndpointCard method="POST" path="/api/v1/auth/register" desc="Register the initial admin user. Fails if a user already exists." auth={false}
            req={`{\n  "username": "admin",\n  "password": "mypassword"\n}`}
            res={`{\n  "message": "User created"\n}`}
            curl={`curl -X POST \\\n  -d '{"username":"admin","password":"mypassword"}' \\\n  http://localhost:8080/api/v1/auth/register`}
          />
          <EndpointCard method="POST" path="/api/v1/auth/login" desc="Authenticate and receive a session token and device info."
            req={`{\n  "username": "admin",\n  "password": "mypassword"\n}`}
            res={`{\n  "token": "a1b2c3d4...",\n  "access_token": "...",\n  "device_info": {\n    "id": "...",\n    "name": "..."\n  }\n}`}
            curl={`curl -X POST \\\n  -d '{"username":"admin","password":"mypassword"}' \\\n  http://localhost:8080/api/v1/auth/login`}
          />
          <EndpointCard method="POST" path="/api/v1/auth/logout" desc="Invalidate the current session token." req={null}
            res={`{\n  "message": "Logged out successfully"\n}`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  http://localhost:8080/api/v1/auth/logout`}
          />
          <EndpointCard method="POST" path="/api/v1/auth/token/generate" desc="Generate a new short-lived access token (5 hours). Requires session token, not access token." req={null}
            res={`{\n  "token": "new_access_token_hex"\n}`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $SESSION_TOKEN" \\\n  http://localhost:8080/api/v1/auth/token/generate`}
          />
          <EndpointCard method="POST" path="/api/v1/auth/token/remove" desc="Remove all access tokens for a given token."
            req={`{\n  "token": "token_to_clean"\n}`}
            res={`{\n  "message": "Access tokens cleaned"\n}`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"token":"token_to_clean"}' \\\n  http://localhost:8080/api/v1/auth/token/remove`}
          />
        </Section>

        {/* ─── Files & Directories ─── */}
        <Section id="files" title="Files & Directories">
          <EndpointCard method="GET" path="/api/v1/files/view?dir={path}" desc="List contents of a directory. Returns files and folders as JSON." req={null}
            res={`[\n  {\n    "name": "Documents",\n    "is_dir": true,\n    "size": 4096,\n    "path": "/home/user",\n    "raw_path": "/home/user/Documents"\n  }\n]`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/view?dir=/home/user"`}
          />
          <EndpointCard method="POST" path="/api/v1/files/directory/create?path={path}" desc="Create a new directory at the specified path." req={null} res={null}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/directory/create?path=/home/user/new_folder"`}
          />
          <EndpointCard method="POST" path="/api/v1/files/unzip" desc="Extract a ZIP archive. If dest_dir is omitted, it extracts to the zip's folder."
            req={`{\n  "zip_path": "/home/user/archive.zip",\n  "dest_dir": "/home/user/output"\n}`}
            res={`{\n  "message": "Archive extracted successfully",\n  "dest": "/home/user/output"\n}`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"zip_path":"/home/user/archive.zip","dest_dir":"/home/user/output"}' \\\n  http://localhost:8080/api/v1/files/unzip`}
          />
        </Section>

        {/* ─── Upload & Download ─── */}
        <Section id="upload-download" title="Upload & Download">
          <EndpointCard method="POST" path="/api/v1/files/upload?dir={path}" desc="Upload a single file via multipart/form-data. Field name: file." req={null} res={null}
            curl={`curl -F "file=@/path/to/file.txt" \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/upload?dir=/home/user"`}
          />
          <EndpointCard method="GET" path="/api/v1/files/download?path={file}" desc="Download a file. Supports token via ?tkn= query param." req={null} res={`<binary file data>`}
            curl={`curl -O \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/download?path=/home/user/file.txt"`}
          />
        </Section>

        {/* ─── Chunked Uploads ─── */}
        <Section id="chunked" title="Chunked Uploads" subtitle="Three-step flow for multi-GB files: initialize, stream parts, then merge.">
          <EndpointCard method="POST" path="/api/v1/files/upload/chunk/start?filename={name}" desc="Step 1 — Initialize a chunked upload session." req={null} res={`Upload initialized`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/upload/chunk/start?filename=big.zip"`}
          />
          <EndpointCard method="POST" path="/api/v1/files/upload/chunk?filename={name}&chunk_index={idx}" desc="Step 2 — Upload a single binary chunk (raw body, NOT multipart). Keep under 90MB." req={null} res={`Chunk received`}
            curl={`curl --data-binary @part0.bin \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/upload/chunk?filename=big.zip&chunk_index=0"`}
          />
          <EndpointCard method="POST" path="/api/v1/files/upload/chunk/complete?filename={name}&dir={path}" desc="Step 3 — Merge all chunks into the final file." req={null} res={`File merged successfully`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/upload/chunk/complete?filename=big.zip&dir=/home/user"`}
          />
        </Section>

        {/* ─── Media & Images ─── */}
        <Section id="media" title="Media & Images">
          <EndpointCard method="GET" path="/api/v1/images/thumbnail?path={file}" desc="Get a resized 200px JPEG thumbnail. Supports ?tkn=<token> for img tags." req={null} res={`<binary image>`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/images/thumbnail?path=/home/user/photo.jpg"`}
          />
          <EndpointCard method="GET" path="/api/v1/images/view?path={file}" desc="Full-resolution image with 7-day client-side cache. Supports ?tkn= query param." req={null} res={`<binary image>`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/images/view?path=/home/user/photo.jpg"`}
          />
          <EndpointCard method="GET" path="/api/v1/files/stream?path={file}" desc="Zero-buffer stream a video or audio file. Supports HTTP Range for seeking and ?tkn= query param." req={null} res={`<binary stream>`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/files/stream?path=/home/user/video.mp4"`}
          />
        </Section>

        {/* ─── Temporary Sharing ─── */}
        <Section id="sharing" title="Temporary Sharing">
          <EndpointCard method="GET" path="/api/v1/share/temp/request?file_path={path}" desc="Generate a one-time temporary download link." auth={false} req={null}
            res={`{\n  "share_link": "/api/v1/share/temp/abcdef...",\n  "token": "abcdef..."\n}`}
            curl={`curl "http://localhost:8080/api/v1/share/temp/request?file_path=/home/user/file.zip"`}
          />
          <EndpointCard method="GET" path="/api/v1/share/temp/{token}" desc="Download via one-time token. No auth required. Token is consumed after use." auth={false} req={null} res={`<binary file>`}
            curl={`curl -O "http://localhost:8080/api/v1/share/temp/abcdef123456"`}
          />
        </Section>

        {/* ─── Anonymous Uploads ─── */}
        <Section id="anonymous" title="Anonymous Uploads" subtitle="Token-based uploads for external users. Admin generates tokens, guests upload with them.">
          <EndpointCard method="POST" path="/api/v1/anonymous/token/generate" desc="Admin: Create a new upload token. Configure via query params."
            req={null}
            res={`{\n  "token": "abc123...",\n  "upload_url": "<host>/upload?token=abc123",\n  "target_dir": "/uploads",\n  "max_uploads": 1,\n  "expiry_hours": 24,\n  "max_file_size": 11811160064\n}`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/anonymous/token/generate?target_dir=/uploads&max_uploads=5&expiry_hours=48&max_file_size=1073741824"`}
          />
          <EndpointCard method="POST" path="/api/v1/anonymous/token/revoke?token={token}" desc="Admin: Revoke an existing upload token." req={null}
            res={`{\n  "message": "Token revoked successfully"\n}`}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/anonymous/token/revoke?token=abc123"`}
          />
          <EndpointCard method="GET" path="/api/v1/anonymous/tokens" desc="Admin: List all active upload tokens." req={null}
            res={`{\n  "tokens": [\n    {\n      "Token": "abc123",\n      "TargetDir": "/uploads",\n      "MaxUploads": 5,\n      "UploadCount": 2,\n      "ExpiresAt": "2026-09-16T...",\n      "IsRevoked": false,\n      "MaxFileSize": 1073741824,\n      "AllowedTypes": []\n    }\n  ],\n  "count": 1\n}`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  http://localhost:8080/api/v1/anonymous/tokens`}
          />
          <EndpointCard method="GET" path="/api/v1/anonymous/token/info?token={token}" desc="Admin: Get details about a specific upload token." req={null}
            res={`{\n  "Token": "abc123",\n  "TargetDir": "/uploads",\n  "MaxUploads": 5,\n  "UploadCount": 2,\n  "ExpiresAt": "2026-09-16T...",\n  "IsRevoked": false,\n  "MaxFileSize": 1073741824,\n  "AllowedTypes": []\n}`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/anonymous/token/info?token=abc123"`}
          />
          <EndpointCard method="POST" path="/api/v1/anonymous/token/validate?token={token}" desc="Guest: Check if a token is valid before uploading." auth={false} req={null}
            res={`{\n  "valid": true,\n  "target_dir": "/uploads",\n  "max_uploads": 5,\n  "uploads_remaining": 3,\n  "expires_at": "2026-09-16T..."\n}`}
            curl={`curl -X POST "http://localhost:8080/api/v1/anonymous/token/validate?token=abc123"`}
          />
          <EndpointCard method="POST" path="/api/v1/anonymous/upload?token={token}" desc="Guest: Upload a file using a valid upload token." auth={false} req={null}
            res={`{\n  "message": "Upload successful",\n  "uploads_remaining": 2\n}`}
            curl={`curl -F "file=@/path/to/file.txt" \\\n  "http://localhost:8080/api/v1/anonymous/upload?token=abc123"`}
          />
          <EndpointCard method="POST" path="/api/v1/anonymous/upload/chunk/start?token={token}&filename={name}" desc="Guest: Start chunked upload with token." auth={false} req={null} res={`Upload initialized`}
            curl={`curl -X POST \\\n  "http://localhost:8080/api/v1/anonymous/upload/chunk/start?token=abc123&filename=big.zip"`}
          />
          <EndpointCard method="POST" path="/api/v1/anonymous/upload/chunk?token={token}&filename={name}&chunk_index={idx}" desc="Guest: Upload a binary chunk." auth={false} req={null} res={`Chunk received`}
            curl={`curl --data-binary @part0.bin \\\n  "http://localhost:8080/api/v1/anonymous/upload/chunk?token=abc123&filename=big.zip&chunk_index=0"`}
          />
          <EndpointCard method="POST" path="/api/v1/anonymous/upload/chunk/complete?token={token}&filename={name}" desc="Guest: Complete chunked upload." auth={false} req={null} res={`File merged successfully`}
            curl={`curl -X POST \\\n  "http://localhost:8080/api/v1/anonymous/upload/chunk/complete?token=abc123&filename=big.zip"`}
          />
        </Section>

        {/* ─── System ─── */}
        <Section id="system" title="System">
          <EndpointCard method="GET" path="/api/v1/system/info" desc="Server version, codename, OS, architecture. No auth required." auth={false} req={null}
            res={`{\n  "version": "v0.0.1-beta",\n  "codename": "Photon",\n  "os": "linux",\n  "arch": "amd64"\n}`}
            curl={`curl http://localhost:8080/api/v1/system/info`}
          />
          <EndpointCard method="GET" path="/api/v1/system/root-path" desc="Get the monitored root directory." req={null}
            res={`{\n  "root_path": "/home/user"\n}`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  http://localhost:8080/api/v1/system/root-path`}
          />
          <EndpointCard method="GET" path="/api/v1/system/storage" desc="Disk usage statistics including external devices." req={null}
            res={`{\n  "total": 512000000000,\n  "used": 256000000000,\n  "available": 256000000000,\n  "usage_percentage": 50.0,\n  "external_devices": [\n    {\n      "device": "/dev/sdb1",\n      "mount_point": "/media/drive",\n      "filesystem_type": "ext4",\n      "total": 64000000000,\n      "used": 32000000000,\n      "available": 32000000000,\n      "usage_percentage": 50.0\n    }\n  ]\n}`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  http://localhost:8080/api/v1/system/storage`}
          />
        </Section>

        {/* ─── Movies ─── */}
        <Section id="movies" title="Movies" subtitle="Legacy endpoints. Not actively developed.">
          <EndpointCard method="GET" path="/api/v1/movies" desc="List all movies." req={null}
            res={`[\n  {\n    "id": 1,\n    "title": "Interstellar",\n    "file_id": 10,\n    "created_at": "2026-01-01T00:00:00Z",\n    "file_path": "/media/movies/interstellar.mp4",\n    "description": "A great movie",\n    "poster_path": "/media/posters/interstellar.jpg",\n    "size": 1073741824,\n    "mime_type": "video/mp4"\n  }\n]`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  http://localhost:8080/api/v1/movies`}
          />
          <EndpointCard method="GET" path="/api/v1/movies/search?q={query}" desc="Search movies by title or description." req={null}
            res={`[\n  {\n    "id": 1,\n    "title": "Interstellar",\n    "file_id": 10,\n    "created_at": "2026-01-01T00:00:00Z",\n    "file_path": "/media/movies/interstellar.mp4",\n    "description": "A great movie",\n    "poster_path": "/media/posters/interstellar.jpg",\n    "size": 1073741824,\n    "mime_type": "video/mp4"\n  }\n]`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/movies/search?q=inter"`}
          />
          <EndpointCard method="POST" path="/api/v1/movies/add" desc="Add a new movie entry."
            req={`{\n  "title": "Interstellar",\n  "file_path": "/media/movies/interstellar.mp4",\n  "file_id": 10,\n  "description": "A great movie",\n  "poster_path": "/media/posters/interstellar.jpg"\n}`}
            res={null}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"title":"Interstellar","file_path":"/media/movies/interstellar.mp4","file_id":10,"description":"A great movie","poster_path":"/media/posters/interstellar.jpg"}' \\\n  http://localhost:8080/api/v1/movies/add`}
          />
          <EndpointCard method="PUT" path="/api/v1/movies/edit" desc="Edit an existing movie."
            req={`{\n  "id": 1,\n  "title": "Updated Title",\n  "description": "Updated description",\n  "poster_path": "/media/posters/updated.jpg"\n}`}
            res={null}
            curl={`curl -X PUT \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"id":1,"title":"Updated Title","description":"Updated description","poster_path":"/media/posters/updated.jpg"}' \\\n  http://localhost:8080/api/v1/movies/edit`}
          />
          <EndpointCard method="GET" path="/api/v1/movies/stream?id={id}" desc="Stream a movie. Supports HTTP Range and ?tkn= query param." req={null} res={`<binary video>`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/movies/stream?id=1"`}
          />
        </Section>

        {/* ─── Series & Episodes ─── */}
        <Section id="series" title="Series & Episodes" subtitle="Legacy endpoints. Not actively developed.">
          <EndpointCard method="GET" path="/api/v1/series" desc="List all series." req={null}
            res={`[\n  {\n    "id": 1,\n    "title": "Breaking Bad",\n    "description": "Crime drama",\n    "created_at": "2026-01-01T00:00:00Z"\n  }\n]`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  http://localhost:8080/api/v1/series`}
          />
          <EndpointCard method="GET" path="/api/v1/series/search?q={query}" desc="Search series by title." req={null}
            res={`[\n  {\n    "id": 1,\n    "title": "Breaking Bad",\n    "description": "Crime drama",\n    "created_at": "2026-01-01T00:00:00Z"\n  }\n]`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/series/search?q=breaking"`}
          />
          <EndpointCard method="POST" path="/api/v1/series/add" desc="Create a new series."
            req={`{\n  "title": "Stranger Things",\n  "description": "Sci-fi horror drama."\n}`} res={null}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"title":"Stranger Things","description":"Sci-fi horror drama."}' \\\n  http://localhost:8080/api/v1/series/add`}
          />
          <EndpointCard method="PUT" path="/api/v1/series/edit" desc="Edit an existing series."
            req={`{\n  "id": 1,\n  "title": "Updated Title",\n  "description": "Updated description"\n}`} res={null}
            curl={`curl -X PUT \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"id":1,"title":"Updated Title","description":"Updated description"}' \\\n  http://localhost:8080/api/v1/series/edit`}
          />
          <EndpointCard method="GET" path="/api/v1/series/episodes?series_id={id}" desc="List episodes for a series." req={null}
            res={`[\n  {\n    "id": 101,\n    "series_id": 1,\n    "file_id": 10,\n    "season_number": 1,\n    "episode_number": 1,\n    "title": "Pilot",\n    "description": "Episode description",\n    "file_path": "/path/to/episode.mp4",\n    "size": 524288000,\n    "mime_type": "video/mp4",\n    "created_at": "2026-01-01T00:00:00Z"\n  }\n]`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/series/episodes?series_id=1"`}
          />
          <EndpointCard method="POST" path="/api/v1/series/episodes/add" desc="Add an episode to a series."
            req={`{\n  "series_id": 1,\n  "file_id": 10,\n  "file_path": "/media/series/s01e01.mp4",\n  "season_number": 1,\n  "episode_number": 1,\n  "title": "Chapter One",\n  "description": "The first chapter"\n}`} res={null}
            curl={`curl -X POST \\\n  -H "Authorization: Bearer $TOKEN" \\\n  -d '{"series_id":1,"file_id":10,"file_path":"/media/series/s01e01.mp4","season_number":1,"episode_number":1,"title":"Chapter One","description":"The first chapter"}' \\\n  http://localhost:8080/api/v1/series/episodes/add`}
          />
          <EndpointCard method="GET" path="/api/v1/series/episodes/stream?id={id}" desc="Stream an episode. Supports HTTP Range and ?tkn= query param." req={null} res={`<binary video>`}
            curl={`curl -H "Authorization: Bearer $TOKEN" \\\n  "http://localhost:8080/api/v1/series/episodes/stream?id=101"`}
          />
        </Section>

        {/* ─── Music ─── */}
        <Section id="music" title="Music" subtitle="Music tracks, streaming, and play statistics. No auth middleware on these routes.">
          <EndpointCard method="GET" path="/api/v1/music" desc="Get all music tracks." auth={false} req={null}
            res={`[\n  {\n    "id": 1,\n    "title": "Bohemian Rhapsody",\n    "artist": "Queen",\n    "album": "...",\n    "duration": 354.5,\n    "releaseYear": 1975,\n    "filePath": "/path/to/song.mp3",\n    "ytUri": "https://youtube.com/watch?v=...",\n    "imagePath": "/path/to/cover.jpg",\n    "size": 8421376,\n    "createdAt": "2026-09-15T..."\n  }\n]`}
            curl={`curl http://localhost:8080/api/v1/music`}
          />
          <EndpointCard method="GET" path="/api/v1/music/get?id={id}" desc="Get a single track by ID." auth={false} req={null}
            res={`{\n  "id": 1,\n  "title": "Bohemian Rhapsody",\n  "artist": "Queen",\n  "album": "...",\n  "duration": 354.5,\n  "releaseYear": 1975,\n  "filePath": "/path/to/song.mp3",\n  "ytUri": "https://youtube.com/watch?v=...",\n  "imagePath": "/path/to/cover.jpg",\n  "size": 8421376,\n  "createdAt": "2026-09-15T..."\n}`}
            curl={`curl "http://localhost:8080/api/v1/music/get?id=1"`}
          />
          <EndpointCard method="GET" path="/api/v1/music/search?title={query}" desc="Search tracks by title. Note: query param is 'title', not 'q'." auth={false} req={null}
            res={`[\n  {\n    "id": 1,\n    "title": "Bohemian Rhapsody",\n    "artist": "Queen",\n    "album": "...",\n    "duration": 354.5,\n    "releaseYear": 1975,\n    "filePath": "/path/to/song.mp3",\n    "ytUri": "https://youtube.com/watch?v=...",\n    "imagePath": "/path/to/cover.jpg",\n    "size": 8421376,\n    "createdAt": "2026-09-15T..."\n  }\n]`}
            curl={`curl "http://localhost:8080/api/v1/music/search?title=bohemian"`}
          />
          <EndpointCard method="POST" path="/api/v1/music/add?uri={youtube_url}" desc="Add a track by URI (downloads via yt-dlp). Uses form value, not JSON." auth={false} req={null}
            res={`{\n  "status": "success",\n  "message": "Queued 1 tracks for download",\n  "count": 1\n}`}
            curl={`curl -X POST \\\n  "http://localhost:8080/api/v1/music/add?uri=https://youtube.com/watch?v=..."`}
          />
          <EndpointCard method="PUT" path="/api/v1/music/update" desc="Update an existing track's metadata via JSON body." auth={false}
            req={`{\n  "id": 1,\n  "title": "Updated Title",\n  "artist": "Updated Artist",\n  "album": "Updated Album",\n  "duration": 210.5,\n  "releaseYear": 2024,\n  "filePath": "/path/to/song.mp3",\n  "ytUri": "https://youtube.com/watch?v=...",\n  "imagePath": "/path/to/cover.jpg",\n  "size": 5000000,\n  "createdAt": "2026-09-15T..."\n}`}
            res={`{\n  "id": 1,\n  "title": "Updated Title",\n  "artist": "Updated Artist",\n  "album": "Updated Album",\n  "duration": 210.5,\n  "releaseYear": 2024,\n  "filePath": "/path/to/song.mp3",\n  "ytUri": "https://youtube.com/watch?v=...",\n  "imagePath": "/path/to/cover.jpg",\n  "size": 5000000,\n  "createdAt": "2026-09-15T..."\n}`}
            curl={`curl -X PUT \\\n  -d '{"id":1,"title":"Updated Title","artist":"Updated Artist","album":"Updated Album","duration":210.5,"releaseYear":2024,"filePath":"/path/to/song.mp3","ytUri":"https://youtube.com/watch?v=...","imagePath":"/path/to/cover.jpg","size":5000000,"createdAt":"2026-09-15T..."}' \\\n  http://localhost:8080/api/v1/music/update`}
          />
          <EndpointCard method="GET" path="/api/v1/music/stream?id={id}" desc="Stream an audio track via http.ServeFile." auth={false} req={null} res={`<binary audio>`}
            curl={`curl "http://localhost:8080/api/v1/music/stream?id=1"`}
          />
          <EndpointCard method="POST" path="/api/v1/music/stats/play?track_id={id}" desc="Record a play event. Params via query or form values." auth={false} req={null}
            res={`{\n  "status": "success",\n  "message": "Play recorded successfully"\n}`}
            curl={`curl -X POST \\\n  "http://localhost:8080/api/v1/music/stats/play?track_id=1"`}
          />
          <EndpointCard method="POST" path="/api/v1/music/stats/favorite?track_id={id}&is_favorite={bool}" desc="Toggle favorite status. Params via query or form values." auth={false} req={null}
            res={`{\n  "status": "success",\n  "message": "Favorite status updated"\n}`}
            curl={`curl -X POST \\\n  "http://localhost:8080/api/v1/music/stats/favorite?track_id=1&is_favorite=true"`}
          />
          <EndpointCard method="GET" path="/api/v1/music/stats/play-state?track_id={id}" desc="Get play state for a track." auth={false} req={null}
            res={`{\n  "track_id": 1,\n  "user_id": 1,\n  "play_count": 12,\n  "last_played_at": "2026-09-15T..."\n}`}
            curl={`curl "http://localhost:8080/api/v1/music/stats/play-state?track_id=1"`}
          />
          <EndpointCard method="GET" path="/api/v1/music/stats/play-states" desc="Get play states for all tracks." auth={false} req={null}
            res={`[\n  {\n    "track_id": 1,\n    "user_id": 1,\n    "play_count": 12,\n    "last_played_at": "2026-09-15T..."\n  }\n]`}
            curl={`curl http://localhost:8080/api/v1/music/stats/play-states`}
          />
          <EndpointCard method="GET" path="/api/v1/music/stats/is-favorite?track_id={id}" desc="Check if a track is favorited." auth={false} req={null}
            res={`{\n  "is_favorite": true\n}`}
            curl={`curl "http://localhost:8080/api/v1/music/stats/is-favorite?track_id=1"`}
          />
          <EndpointCard method="GET" path="/api/v1/music/stats/last-played" desc="Get the last played track." auth={false} req={null}
            res={`{\n  "user_id": 1,\n  "track_id": 3,\n  "played_at": "2026-09-15T..."\n}`}
            curl={`curl http://localhost:8080/api/v1/music/stats/last-played`}
          />
        </Section>

        {/* ─── WebSockets & Devices ─── */}
        <Section id="websockets" title="WebSockets & Devices">
          <EndpointCard method="GET" path="/api/v1/devices" desc="List all currently connected WebSocket devices." auth={false} req={null}
            res={`[\n  {\n    "deviceId": "abc-123",\n    "deviceName": "Swap's Phone"\n  }\n]`}
            curl={`curl http://localhost:8080/api/v1/devices`}
          />
          <EndpointCard method="GET" path="/api/v1/ws?deviceId={id}&deviceName={name}" desc="WebSocket endpoint. Connect with optional deviceId and deviceName query params." auth={false} req={null}
            res={`101 Switching Protocols (WebSocket)`}
            curl={`wscat -c "ws://localhost:8080/api/v1/ws?deviceId=phone-1&deviceName=MyPhone"`}
          />
        </Section>

      </main>

      <footer className="border-t border-white/5 py-10 text-center relative z-10">
        <p className="text-gray-600 text-xs font-medium">&copy; 2026 CodersHub Inc. Open Source under the AGPLv3 License.</p>
      </footer>
    </div>
  );
}

function NavLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <a href={href} className="px-2.5 py-1.5 rounded-md text-[11px] font-medium text-gray-500 hover:text-white hover:bg-white/5 transition-all">
      {children}
    </a>
  );
}

function Section({ id, title, subtitle, children }: { id: string; title: string; subtitle?: string; children: React.ReactNode }) {
  return (
    <section id={id} className="mb-20 scroll-mt-20">
      <h2 className="text-xl font-bold mb-1 text-white tracking-tight">{title}</h2>
      {subtitle && <p className="text-gray-500 text-sm mb-1">{subtitle}</p>}
      <Separator className="mb-6 bg-white/[0.06]" />
      <div className="space-y-4">{children}</div>
    </section>
  );
}

function EndpointCard({ method, path, desc, req, res, curl, auth = true }: { method: string; path: string; desc: string; req: string | null; res: string | null; curl: string; auth?: boolean }) {
  const badge: Record<string, string> = {
    GET:    'bg-sky-500/10 text-sky-400 border-sky-500/20',
    POST:   'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
    PUT:    'bg-amber-500/10 text-amber-400 border-amber-500/20',
    DELETE: 'bg-rose-500/10 text-rose-400 border-rose-500/20',
  };

  return (
    <div className="bg-white/[0.015] border border-white/[0.05] rounded-xl overflow-hidden hover:border-white/10 transition-all">
      <div className="px-4 py-3 flex flex-col sm:flex-row sm:items-center gap-3">
        <div className="flex items-center gap-2.5 flex-1 min-w-0 flex-wrap">
          <span className={`shrink-0 text-[10px] font-bold px-2 py-0.5 rounded border ${badge[method] || ''}`}>{method}</span>
          <code className="text-indigo-300/90 font-mono text-xs break-all leading-relaxed">{path}</code>
          {!auth && <span className="shrink-0 text-[9px] font-medium px-1.5 py-0.5 rounded bg-yellow-500/10 text-yellow-500 border border-yellow-500/20">PUBLIC</span>}
        </div>
      </div>

      <div className="px-4 pb-3 pt-0">
        <p className="text-xs text-gray-500 mb-3">{desc}</p>
        <div className="flex items-start gap-2">
          <pre className="flex-1 bg-black/40 text-gray-400 px-3 py-2.5 rounded-lg text-[11px] overflow-x-auto font-mono border border-white/[0.04] leading-relaxed">
            <code>{curl}</code>
          </pre>
          <div className="shrink-0 pt-1">
            <CopyButton text={curl.replace(/\\\n\s*/g, '')} />
          </div>
        </div>

        {(req || res) && (
          <div className="flex flex-wrap gap-2 mt-3">
            {req && <ExpandableCode label="Request Body" code={req} variant="request" />}
            {res && <ExpandableCode label="Response" code={res} variant="response" />}
          </div>
        )}
      </div>
    </div>
  );
}
