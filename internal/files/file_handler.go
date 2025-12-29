package files

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Name    string `json:"name"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size"`
	Path    string `json:"path"`
	RawPath string `json:"raw_path"`
}

func viewContent(dir string) ([]FileInfo, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		file := FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			Path:    dir,
			RawPath: filepath.Join(dir, entry.Name()),
		}
		files = append(files, file)
	}
	return files, nil
}

func createDir(path string) error {
	const mode = 0755

	err := os.MkdirAll(path, mode)
	if err != nil {
		log.Printf("Failed to create directory '%s': %v", path, err)
		return fmt.Errorf("create directory failed: %w", err)
	}

	log.Printf("Directory ready: %s", path)
	return nil
}

// UploadLargeFileToDir handles massive files (10GB+) efficiently
func uploadLargeFileToDir(w http.ResponseWriter, r *http.Request, targetDir string) error {
	// 1. Limit the size of incoming request to 11GB
	const maxUploadSize = 11 << 30

	// This protects the server from accepting >11GB, but strictly streams it.
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// 2. Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target dir: %w", err)
	}

	// 3. Get the Multipart Reader (Streaming Mode)
	reader, err := r.MultipartReader()
	if err != nil {
		return fmt.Errorf("failed to start stream (check Content-Type): %w", err)
	}

	// 4. Stream Loop
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("streaming error: %w", err)
		}

		filename := part.FileName()
		if filename == "" {
			continue
		}
		// 6. Direct Stream to Disk
		dstPath := filepath.Join(targetDir, filename)

		// Open the file for writing
		dst, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("failed to create file on disk: %w", err)
		}

		// This is the critical line for 11GB files.
		// It copies 32KB chunks at a time from Network -> Disk.
		// RAM usage remains near zero.
		if _, err := io.Copy(dst, part); err != nil {
			dst.Close()
			return fmt.Errorf("upload interrupted: %w", err)
		}
		dst.Close()
	}

	return nil
}
