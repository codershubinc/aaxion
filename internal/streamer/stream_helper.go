package streamer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Define a chunk size (e.g., 1MB) to prevent sending the whole file at once
const MAX_CHUNK_SIZE = 1 * 1024 * 1024

func isBrowserNative(filename string) bool {
	browserNativeFormats := []string{
		".mp4",
		".webm",
		".ogg",
		".mp3",
		".wav",
		".aac",
		".flac",
		".m4a",
		".mov",
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return slices.Contains(browserNativeFormats, ext)
}

// StreamFileRange handles range requests manually
func StreamFileRange(w http.ResponseWriter, r *http.Request, filePath string, mimeType string) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not stats file", http.StatusInternalServerError)
		return
	}
	totalSize := stat.Size()

	rangeHeader := r.Header.Get("Range")

	if rangeHeader == "" {
		// Optional: You might want to default to sending the first chunk
		// instead of an error if you want to support direct downloads/playback without headers,
		// but keeping your current logic is fine too.
		log.Println("No range headers")
		http.Error(w, "Range header required", http.StatusBadRequest)
		return
	}
	// log.Println("Range headers found", rangeHeader)

	re := regexp.MustCompile(`bytes=(\d+)-(\d*)`)
	matches := re.FindStringSubmatch(rangeHeader)
	if matches == nil {
		http.Error(w, "Invalid Range header", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	start, _ := strconv.ParseInt(matches[1], 10, 64)
	var end int64

	// --- FIX START ---
	if matches[2] != "" {
		// Client requested specific end
		parsedEnd, _ := strconv.ParseInt(matches[2], 10, 64)
		end = parsedEnd
	} else {
		// Client requested "rest of file" (bytes=0-)
		// We MUST cap this, or we send the whole 2GB file
		end = start + MAX_CHUNK_SIZE - 1
	}

	// Safety check: ensure we don't go past the actual file size
	if end >= totalSize {
		end = totalSize - 1
	}
	// --- FIX END ---

	if start > end {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", totalSize))
		http.Error(w, "Invalid Range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	contentLength := end - start + 1

	// Seek to the start position
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		http.Error(w, "Seek error", http.StatusInternalServerError)
		return
	}

	// Headers
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusPartialContent) // 206

	// Copy ONLY the specific chunk size
	io.CopyN(w, file, contentLength)
}
