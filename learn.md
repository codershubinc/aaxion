# 📚 Learn Aaxion

Welcome to the Aaxion learning guide! This document will help you understand Aaxion, its architecture, and how to effectively use it to turn old hardware into a powerful file server.

---

## 🎯 What is Aaxion?

**Aaxion** is a lightweight, high-performance file server built in Go that transforms old hardware (like laptops, old PCs) into efficient storage nodes. It's designed with extreme resource efficiency in mind - capable of handling 10GB file transfers while using only ~32KB of RAM through zero-buffer streaming.

### Key Benefits

- **Repurpose Old Hardware:** Give new life to old laptops and computers by turning them into dedicated file storage
- **Minimal Resource Usage:** Runs efficiently with ~10MB RAM idle, perfect for low-spec hardware
- **Modern Features:** Supports chunked uploads, resumable transfers, and temporary sharing
- **Cross-Platform:** Works on Linux (primary) and Windows (experimental)

---

## 🏗️ Architecture Overview

### Core Components

1. **HTTP Server (Port 8080)**
   - RESTful API for file operations
   - Streaming-based transfers (no buffering)
   - Built-in security (path sanitization, hidden file exclusion)

2. **File System Monitor**
   - Watches a specified root directory (e.g., `/home/swap/`)
   - Serves file tree via JSON APIs
   - Automatic path validation

3. **Database (SQLite)**
   - Stores temporary share tokens
   - Manages chunked upload sessions
   - Lightweight and embedded

4. **Image Processing**
   - On-demand thumbnail generation
   - Server-side caching for performance
   - Automatic content-type detection

### How Streaming Works

Aaxion uses Go's `io.Copy` to stream data directly between network and disk without loading into memory:

```
Client → Network → Disk (Upload)
Disk → Network → Client (Download)
```

This allows handling files of any size with minimal RAM usage.

---

## 🚀 Getting Started

### Step 1: Download and Install

1. Go to [Releases](https://github.com/codershubinc/aaxion/releases)
2. Download the binary for your operating system:
   - Linux: `aaxion-linux-amd64`
   - Windows: `aaxion-windows-amd64.exe`

### Step 2: Set Permissions (Linux only)

```bash
chmod +x aaxion-linux-amd64
```

### Step 3: Run the Server

```bash
# Linux
./aaxion-linux-amd64

# Windows
./aaxion-windows-amd64.exe
```

The server will start on `http://localhost:8080`

### Step 4: Test the Installation

```bash
# Check the root path
curl http://localhost:8080/api/system/get-root-path

# View files in a directory
curl "http://localhost:8080/files/view?dir=/home/swap/"
```

---

## 📖 Common Use Cases

### Use Case 1: Personal Cloud Storage

Transform your old laptop into a personal cloud:

1. Install Aaxion on the old laptop
2. Set it up as a systemd service (Linux) or startup program (Windows)
3. Connect via the [Aaxion Mobile App](https://github.com/codershubinc/aaxion-mob)
4. Access your files from anywhere on your local network

### Use Case 2: Media Server

Store and stream your media collection:

```bash
# Upload a video
curl -F "file=@movie.mp4" "http://localhost:8080/files/upload?dir=/home/swap/media"

# Generate a temporary share link
curl "http://localhost:8080/files/d/r?file_path=/home/swap/media/movie.mp4"
```

### Use Case 3: Large File Transfers

Transfer large files reliably using chunked uploads:

1. Start upload session
2. Upload chunks (up to 90MB each)
3. Complete the upload to merge chunks

See [API Documentation](./docs/api.md) for detailed examples.

---

## 🔧 Configuration

### Setting the Root Directory

By default, Aaxion monitors a specific directory. You can configure this when starting the server or through environment variables (check the source code for specific configuration options).

**Example configurations:**
```bash
# Set via environment variable (if supported)
export AAXION_ROOT_PATH=/home/swap
./aaxion-linux-amd64

# Or check the command-line flags
./aaxion-linux-amd64 --help
```

### Running as a Service (Linux)

For persistent operation, set up Aaxion as a systemd service:

1. Create a service file at `/etc/systemd/system/aaxion.service`:

```ini
[Unit]
Description=Aaxion File Server
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/path/to/aaxion
ExecStart=/path/to/aaxion/aaxion-linux-amd64
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

2. Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable aaxion
sudo systemctl start aaxion
sudo systemctl status aaxion
```

3. Aaxion will now run in the background and start automatically on boot

---

## 🎓 Tutorials

### Tutorial 1: Basic File Operations

**Upload a file:**
```bash
curl -F "file=@document.pdf" "http://localhost:8080/files/upload?dir=/home/swap/documents"
```

**Create a directory:**
```bash
curl -X POST "http://localhost:8080/files/create-directory?path=/home/swap/new_folder"
```

**List files:**
```bash
curl "http://localhost:8080/files/view?dir=/home/swap/documents"
```

### Tutorial 2: Working with Large Files

**For files > 90MB, use chunked upload:**

```bash
# 1. Split your file
split -b 90M largefile.zip chunk_

# 2. Start upload session
curl -X POST "http://localhost:8080/files/upload/chunk/start?filename=largefile.zip"

# 3. Upload chunks
curl --data-binary @chunk_aa "http://localhost:8080/files/upload/chunk?filename=largefile.zip&chunk_index=0"
curl --data-binary @chunk_ab "http://localhost:8080/files/upload/chunk?filename=largefile.zip&chunk_index=1"

# 4. Complete upload
curl -X POST "http://localhost:8080/files/upload/chunk/complete?filename=largefile.zip&dir=/home/swap/uploads"
```

### Tutorial 3: Image Handling

**View full image:**
```bash
curl "http://localhost:8080/files/view-image?path=/home/swap/photos/vacation.jpg" -o image.jpg
```

**Get thumbnail (200px):**
```bash
curl "http://localhost:8080/files/thumbnail?path=/home/swap/photos/vacation.jpg" -o thumbnail.jpg
```

### Tutorial 4: Temporary File Sharing

**Generate a one-time share link:**
```bash
# Request a token
curl "http://localhost:8080/files/d/r?file_path=/home/swap/document.pdf"

# Response: {"share_link": "/files/d/t/TOKEN", "token": "TOKEN"}

# Share this URL with others
http://localhost:8080/files/d/t/TOKEN
```

---

## 🔒 Security Best Practices

1. **Path Validation:** Aaxion automatically prevents directory traversal attacks
2. **Hidden Files:** Files starting with `.` are automatically excluded
3. **Root Restriction:** All operations are confined to the monitored root directory
4. **Network Security:** Consider using a reverse proxy (Nginx/Apache) for HTTPS
5. **Firewall:** Configure firewall rules to restrict access to trusted networks

---

## 🐛 Troubleshooting

### Server won't start
- Check if port 8080 is already in use
- Verify file permissions on the binary
- Check system logs for error messages

### Cannot upload files
- Ensure the target directory exists
- Verify filesystem permissions
- Check available disk space

### Path errors
- All paths must be within the monitored root directory
- Use absolute paths starting with the root (e.g., `/home/swap/`)
- Avoid relative paths like `../`

---

## 📚 Additional Resources

- **API Reference:** [docs/api.md](./docs/api.md) - Complete API documentation
- **Mobile App:** [Aaxion-Mob](https://github.com/codershubinc/aaxion-mob) - Android/iOS client
- **Source Code:** [GitHub Repository](https://github.com/codershubinc/aaxion)
- **Issues & Support:** [GitHub Issues](https://github.com/codershubinc/aaxion/issues)

---

## 💡 Tips and Tricks

1. **Performance:** Use an SSD for better performance, even an old one
2. **Networking:** Connect via Ethernet for stable transfers
3. **Monitoring:** Keep an eye on disk space usage
4. **Backups:** Regular backups are recommended for important data
5. **Updates:** Check for new releases regularly for bug fixes and features

---

## 🤝 Contributing

Want to improve Aaxion? Contributions are welcome! Check the repository for contribution guidelines.

---

**Ready to get started?** Install Aaxion and transform your old hardware into a powerful file server today! 🚀
