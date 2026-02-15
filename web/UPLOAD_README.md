# Aaxion Public File Upload Interface

A modern, real-time web interface for accepting file uploads from unknown users (no authentication required) with speed monitoring and progress tracking.

## Features

### 📤 Public File Upload

- **No Authentication Required** - Accept files from anyone
- Drag and drop or click to browse
- Multiple file uploads simultaneously
- Supports files up to 11GB
- Real-time progress tracking

### ⚡ Smart Upload Mode

- **Direct Upload**: For files under 100MB (faster, single request)
- **Chunked Upload**: For files over 100MB (reliable, resumable)
- Automatic mode switching based on file size
- Manual mode toggle available

### 📊 Real-Time Statistics

- **Upload Speed**: Live MB/s per file
- **Progress Tracking**: Visual progress bars and percentages
- **Total Statistics**: Overall upload speed and data transferred
- **Queue Management**: See all uploads in progress

### 🎨 Modern UI

- Beautiful gradient design with Tailwind CSS
- Drag and drop support
- Toast notifications
- Responsive layout
- Real-time updates

## API Endpoints

### Direct Upload (No Authentication)

```
POST /files/upload/public?dir={directory}
Content-Type: multipart/form-data

Form field: file
```

### Chunked Upload (No Authentication)

#### 1. Initialize Upload

```
POST /files/upload/chunk/start/public?filename={filename}
```

#### 2. Upload Chunks

```
POST /files/upload/chunk/public?filename={filename}&chunk_index={index}
Content-Type: application/octet-stream

Body: raw chunk data
```

#### 3. Complete Upload

```
POST /files/upload/chunk/complete/public?filename={filename}&dir={directory}
```

## Usage

### For Users

1. **Access the upload page**
   Navigate to: `http://your-server:8080/upload`

2. **Select upload directory**
   Enter the server directory where files should be saved (e.g., `/uploads`)

3. **Upload files**
   - Drag and drop files onto the upload zone
   - OR click the upload zone to browse files
   - OR click "Select Files" button

4. **Choose upload mode** (optional)
   - Enable "Chunked Upload" for large files or unreliable connections
   - Leave disabled for automatic mode selection

5. **Monitor progress**
   - Watch real-time upload speed
   - Track progress for each file
   - View overall statistics

### For Developers

#### Integrate the Upload API

**Direct Upload Example (JavaScript):**

```javascript
const formData = new FormData();
formData.append("file", fileObject);

const response = await fetch("/files/upload/public?dir=/uploads", {
  method: "POST",
  body: formData,
});

if (response.ok) {
  console.log("Upload successful!");
}
```

**Chunked Upload Example:**

```javascript
const CHUNK_SIZE = 50 * 1024 * 1024; // 50MB

// 1. Initialize
await fetch(`/files/upload/chunk/start/public?filename=${filename}`, {
  method: "POST",
});

// 2. Upload chunks
const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
for (let i = 0; i < totalChunks; i++) {
  const start = i * CHUNK_SIZE;
  const end = Math.min(start + CHUNK_SIZE, file.size);
  const chunk = file.slice(start, end);

  await fetch(
    `/files/upload/chunk/public?filename=${filename}&chunk_index=${i}`,
    {
      method: "POST",
      body: chunk,
    },
  );
}

// 3. Complete
await fetch(
  `/files/upload/chunk/complete/public?filename=${filename}&dir=/uploads`,
  {
    method: "POST",
  },
);
```

## Configuration

### Upload Directory

The default upload directory is `/uploads`. You can customize this by:

1. Changing the value in the "Upload Directory" field
2. Using query parameter: `?dir=/your/custom/path`

**Security Note**: The server validates paths to prevent directory traversal attacks.

### File Size Limits

- **Direct Upload**: Up to 11GB
- **Chunked Upload**: Unlimited (50MB chunks)

### Chunk Size

Default: 50MB per chunk. Modify in `upload.js`:

```javascript
const CHUNK_SIZE = 50 * 1024 * 1024; // 50MB
```

## File Structure

```
web/
├── upload.html      # Main upload interface
├── upload.js        # Upload functionality with progress tracking
├── styles.css       # Custom CSS styles
└── UPLOAD_README.md # This file
```

## Security Considerations

### ⚠️ Important Security Notes

This interface accepts uploads **without authentication**. Consider these security measures:

1. **Rate Limiting**: Implement rate limiting to prevent abuse
2. **File Type Validation**: Add server-side file type checks
3. **Virus Scanning**: Scan uploaded files for malware
4. **Disk Quotas**: Set upload directory quotas
5. **IP Whitelisting**: Restrict access to trusted IPs if needed
6. **Path Validation**: Server validates paths (already implemented)

### Current Security Features

- ✅ Path traversal protection (`ExpelDotPath`)
- ✅ File size limits (11GB for direct uploads)
- ✅ Chunk size limits (90MB per chunk)

### Recommended Additional Security

Add to your server code:

```go
// Example: File type validation
func validateFileType(filename string) bool {
    allowed := []string{".jpg", ".png", ".pdf", ".mp4"}
    ext := filepath.Ext(filename)
    for _, a := range allowed {
        if ext == a {
            return true
        }
    }
    return false
}

// Example: Rate limiting
var uploadLimiter = rate.NewLimiter(5, 10) // 5 uploads per second, burst of 10

func PublicUploadFile(w http.ResponseWriter, r *http.Request) {
    if !uploadLimiter.Allow() {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    // ... rest of upload logic
}
```

## Browser Compatibility

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

Requires:

- Fetch API
- FormData API
- Drag and Drop API
- ES6+ JavaScript

## Troubleshooting

### Upload fails immediately

- Check if the upload directory exists on the server
- Verify the directory path is correct
- Check server logs for errors

### Upload stops mid-way

- Network interruption - try chunked upload mode
- File too large for direct upload - enable chunked mode
- Check server disk space

### Slow upload speeds

- Check network connection
- Try chunked upload for better reliability
- Check server CPU/disk usage

### "Suspicious path detected" error

- Don't use `..` in directory paths
- Use absolute paths starting with `/`
- Don't include hidden files (starting with `.`)

## Performance Tips

### For Large Files

- ✅ Enable chunked upload mode
- ✅ Use wired connection instead of WiFi
- ✅ Close other bandwidth-intensive applications

### For Multiple Files

- The interface uploads files in parallel
- Each file shows individual progress
- Total statistics show combined metrics

### For Servers

- Ensure sufficient disk space
- Monitor disk I/O performance
- Consider using SSD for upload directories
- Implement cleanup jobs for old uploads

## Customization

### Change Upload Chunk Size

Edit in [web/upload.js](web/upload.js):

```javascript
const CHUNK_SIZE = 100 * 1024 * 1024; // 100MB chunks
```

### Change Chunked Threshold

```javascript
const CHUNKED_THRESHOLD = 200 * 1024 * 1024; // Switch to chunked at 200MB
```

### Customize UI Colors

The interface uses Tailwind CSS. Modify colors in the HTML or add custom CSS.

### Add File Type Restrictions (Client-side)

```html
<input type="file" id="fileInput" multiple accept=".jpg,.png,.pdf" />
```

## Future Enhancements

- [ ] Pause/resume uploads
- [ ] Upload cancellation
- [ ] File type filtering
- [ ] Maximum file size configuration
- [ ] Upload history
- [ ] Email notifications on upload
- [ ] Webhook support for integrations
- [ ] Password-protected uploads
- [ ] Expiring upload links
- [ ] Upload quota per IP

## Testing

### Test Direct Upload

1. Select a file under 100MB
2. Disable "Chunked Upload"
3. Upload and verify speed tracking

### Test Chunked Upload

1. Select a file over 100MB
2. Enable "Chunked Upload"
3. Upload and verify chunk progress

### Test Multiple Files

1. Select multiple files
2. Watch parallel upload progress
3. Verify all files complete successfully

## Support

For issues or questions:

- Check server logs for detailed error messages
- Verify network connectivity
- Ensure sufficient disk space
- Review security settings

## License

Part of the Aaxion project.
