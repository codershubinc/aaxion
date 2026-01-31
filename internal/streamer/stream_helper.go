package streamer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

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
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", strconv.FormatInt(totalSize, 10))
		w.Header().Set("Accept-Ranges", "bytes")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, file)
		return
	}

	re := regexp.MustCompile(`bytes=(\d+)-(\d*)`)
	matches := re.FindStringSubmatch(rangeHeader)
	if matches == nil {
		http.Error(w, "Invalid Range header", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	start, _ := strconv.ParseInt(matches[1], 10, 64)
	end := int64(-1)
	if matches[2] != "" {
		end, _ = strconv.ParseInt(matches[2], 10, 64)
	}

	if end == -1 || end >= totalSize {
		end = totalSize - 1
	}

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

	// Copy specific amount of bytes
	io.CopyN(w, file, contentLength)
}
