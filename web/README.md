# Aaxion Web Interfaces

Modern web interfaces for testing and using Aaxion's file management features with real-time monitoring and beautiful UI.

## 🌐 Available Interfaces

### 1. Landing Page (`/`)

**File**: [landing.html](landing.html)

Overview page with links to all web interfaces and system information.

- **URL**: `http://localhost:8080/`
- **Features**:
  - Quick access to all interfaces
  - System features overview
  - API endpoint information

### 2. Public Upload Interface (`/upload`)

**Files**: [upload.html](upload.html), [upload.js](upload.js)

Accept file uploads from unknown users without authentication.

- **URL**: `http://localhost:8080/upload`
- **Features**:
  - 📤 Upload without authentication
  - 📊 Real-time speed monitoring
  - 🔄 Automatic chunked upload for large files
  - 📁 Multiple file uploads
  - 🎯 Drag and drop support
  - 📈 Progress tracking per file
  - ⚡ Supports files up to 11GB

**Documentation**: [UPLOAD_README.md](UPLOAD_README.md)

### 3. File Share Testing (`/web`)

**Files**: [index.html](index.html), [app.js](app.js)

Browse files, generate temporary share links, and test downloads.

- **URL**: `http://localhost:8080/web`
- **Features**:
  - 🗂️ File browser with navigation
  - 🔗 Generate temporary share links (single-use tokens)
  - 📥 Test downloads with speed monitoring
  - 📊 Real-time download statistics
  - 🔒 Token-based secure sharing
  - 📋 One-click copy to clipboard

## 🚀 Quick Start

1. **Start the Aaxion server**

   ```bash
   cd /home/swap/Github/aaxion
   go run cmd/main.go
   ```

2. **Access the landing page**

   ```
   http://localhost:8080/
   ```

3. **Choose your interface**
   - Click "Public Upload" for file uploads
   - Click "File Share" for testing temp file sharing

## 📁 File Structure

```
web/
├── landing.html        # Main landing page
├── upload.html         # Public upload interface
├── upload.js          # Upload functionality
├── index.html         # File share testing interface
├── app.js            # File share functionality
├── styles.css         # Shared custom styles
├── README.md          # This file
└── UPLOAD_README.md   # Detailed upload documentation
```

## 🎨 Tech Stack

- **HTML5**: Semantic markup
- **Tailwind CSS**: Utility-first CSS framework (via CDN)
- **Vanilla JavaScript**: No frameworks, pure JS
- **Font Awesome**: Icons (via CDN)
- **Fetch API**: Modern HTTP requests
- **Streams API**: Real-time monitoring

## 🔌 API Endpoints Used

### Public Upload (No Auth)

```
POST /files/upload/public?dir={directory}
POST /files/upload/chunk/start/public?filename={filename}
POST /files/upload/chunk/public?filename={filename}&chunk_index={index}
POST /files/upload/chunk/complete/public?filename={filename}&dir={directory}
```

### File Share (Requires Auth for generation)

```
GET  /api/files/view?dir={path}           # Browse files (auth required)
GET  /files/d/r?file_path={path}          # Generate share link (auth required)
GET  /files/d/t/{token}                   # Download with token (no auth)
```

## 📊 Features Breakdown

### Public Upload Interface

#### Direct Upload Mode

- For files under 100MB
- Single HTTP request
- Fastest for small files
- Automatic retry on failure

#### Chunked Upload Mode

- For files over 100MB
- Splits file into 50MB chunks
- Reliable for large files
- Can be resumed on failure

#### Real-Time Statistics

- Upload speed per file (MB/s)
- Individual progress bars
- Total uploaded data
- Files in queue counter
- Completed uploads counter

### File Share Interface

#### File Browser

- Navigate directory structure
- View file sizes and types
- Quick select for sharing
- Refresh directory contents

#### Share Link Generation

- Generate temporary tokens
- Single-use download links
- Automatic link expiration
- Copy to clipboard

#### Download Testing

- Real-time speed monitoring
- Progress tracking
- Size and time statistics
- Automatic file download

## 🎨 UI Features

- **Dark Theme**: Optimized for extended use
- **Gradient Backgrounds**: Modern aesthetic
- **Smooth Animations**: Polished transitions
- **Toast Notifications**: User feedback
- **Responsive Design**: Works on all devices
- **Accessibility**: Keyboard navigation support

## 🔒 Security Notes

### Public Upload Interface

⚠️ **Warning**: Accepts uploads without authentication

**Security measures**:

- Path traversal protection (`ExpelDotPath`)
- File size limits (11GB max)
- Chunk size limits (90MB max)

**Recommended additions**:

- Rate limiting per IP
- File type validation
- Virus scanning
- Disk quotas
- IP whitelisting
- CAPTCHA for bot prevention

### File Share Interface

- Generate links requires authentication
- Single-use tokens (auto-revoked after download)
- Server-side validation
- Path traversal protection

## 🌐 Browser Compatibility

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

**Requirements**:

- Fetch API
- Streams API
- FormData API
- Drag and Drop API
- ES6+ JavaScript
- Async/Await

## 🛠️ Development

### Adding New Interfaces

1. **Create HTML file** in `web/` directory
2. **Create JavaScript file** for functionality
3. **Add route** in [internal/api/routes.go](../internal/api/routes.go)
4. **Update landing page** with link to new interface

### Modifying Existing Interfaces

1. **Edit HTML** for structure changes
2. **Edit JavaScript** for functionality changes
3. **Edit styles.css** for custom styling
4. **Test thoroughly** before deploying

### Using Shared Styles

Import the shared styles in your HTML:

```html
<link rel="stylesheet" href="styles.css" />
```

Available custom classes:

- `.scrollbar-thin` - Custom scrollbar
- `card-hover` - Card hover effects
- Animation utilities

## 📝 Configuration

### Upload Settings

**Default upload directory**: `/uploads`

Change in [upload.html](upload.html):

```html
<input type="text" id="uploadDir" value="/uploads" />
```

**Chunk size**: 50MB

Change in [upload.js](upload.js):

```javascript
const CHUNK_SIZE = 50 * 1024 * 1024; // 50MB
```

**Chunked threshold**: 100MB

Change in [upload.js](upload.js):

```javascript
const CHUNKED_THRESHOLD = 100 * 1024 * 1024; // 100MB
```

### File Share Settings

**Default browse directory**: `/`

Auto-loads on page load. Users can navigate to any directory.

## 🐛 Troubleshooting

### Styles not loading

- Check that `styles.css` route is registered
- Verify the file path is correct
- Clear browser cache

### JavaScript not working

- Check browser console for errors
- Verify `.js` route is registered with correct MIME type
- Ensure server is serving static files

### Upload fails

- Check server disk space
- Verify upload directory permissions
- Check server logs for errors
- Try chunked upload mode

### Share links don't work

- Tokens are single-use only
- Generate new link after use
- Check authentication for link generation

## 📈 Performance Tips

### For Users

- Use wired connection for large uploads
- Close bandwidth-intensive applications
- Enable chunked upload for files > 100MB

### For Developers

- Minimize JavaScript bundle size
- Use CDN for libraries
- Enable gzip compression on server
- Optimize images and assets
- Use browser caching

## 🚧 Future Enhancements

- [ ] Pause/resume uploads
- [ ] Upload history and management
- [ ] File preview before upload
- [ ] Folder upload support
- [ ] Link expiration customization
- [ ] QR code for share links
- [ ] Email notifications
- [ ] WebSocket for real-time updates
- [ ] Dark/light theme toggle
- [ ] Internationalization (i18n)

## 📚 Documentation

- **Main README**: [../README.md](../README.md)
- **Upload Guide**: [UPLOAD_README.md](UPLOAD_README.md)
- **API Docs**: [../docs/api.md](../docs/api.md)

## 🤝 Contributing

To add new features to the web interfaces:

1. Follow the existing code style
2. Use Tailwind CSS for styling
3. Add proper error handling
4. Include loading states
5. Update documentation
6. Test on multiple browsers

## 📄 License

Part of the Aaxion project.
