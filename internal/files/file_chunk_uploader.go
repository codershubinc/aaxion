package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	TempDir      = "./uploads_temp"
	MaxChunkSize = 90 << 20
)

func HandleStartChunkUpload(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Create a unique temporary folder for this specific file's chunks
	// We use the filename as the folder name (in production, use a UUID)
	tempFolderPath := filepath.Join(TempDir, filename)
	if err := os.MkdirAll(tempFolderPath, 0755); err != nil {
		http.Error(w, "Failed to create temp storage", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Upload initialized"))
}

func HandleUploadChunk(w http.ResponseWriter, r *http.Request) {
	// Parse query params: ?filename=video.mp4&chunk_index=0
	filename := r.URL.Query().Get("filename")
	indexStr := r.URL.Query().Get("chunk_index")

	//parsing the index
	chunkIndex, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid chunk index", http.StatusBadRequest)
		return
	}

	//limiter
	r.Body = http.MaxBytesReader(w, r.Body, MaxChunkSize)

	// Save the chunk: ./uploads_temp/video.mp4/chunk_0
	tempFolderPath := filepath.Join(TempDir, filename)
	chunkPath := filepath.Join(tempFolderPath, fmt.Sprintf("chunk_%d", chunkIndex))

	dst, err := os.Create(chunkPath)
	if err != nil {
		http.Error(w, "Failed to create chunk file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r.Body); err != nil {
		http.Error(w, "Failed to write chunk", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Chunk received"))
}

// 3. COMPLETE: Merge all chunks
func HandleCompleteUpload(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	uploadDir := r.URL.Query().Get("dir")
	if filename == "" || uploadDir == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	tempFolderPath := filepath.Join(TempDir, filename)
	finalPath := filepath.Join(uploadDir, filename)

	// Ensure final directory exists
	os.MkdirAll(uploadDir, 0755)

	// Create final file
	finalFile, err := os.Create(finalPath)
	if err != nil {
		http.Error(w, "Failed to create final file", http.StatusInternalServerError)
		return
	}
	defer finalFile.Close()

	// 1. List all chunks
	files, err := os.ReadDir(tempFolderPath)
	if err != nil {
		http.Error(w, "Failed to read temp dir", http.StatusInternalServerError)
		return
	}

	// 2. Sort chunks numerically (chunk_1, chunk_2, chunk_10...)
	// Standard string sort is bad because "chunk_10" comes before "chunk_2"
	// So we sort by the index number we parsed manually.
	type ChunkFile struct {
		Path  string
		Index int
	}
	var chunks []ChunkFile

	for _, f := range files {
		var idx int
		// Extract number from "chunk_0", "chunk_1"
		fmt.Sscanf(f.Name(), "chunk_%d", &idx)
		chunks = append(chunks, ChunkFile{
			Path:  filepath.Join(tempFolderPath, f.Name()),
			Index: idx,
		})
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Index < chunks[j].Index
	})

	// 3. Append chunks to final file
	for _, chunk := range chunks {
		// Open chunk
		chunkFile, err := os.Open(chunk.Path)
		if err != nil {
			http.Error(w, "Missing chunk", http.StatusInternalServerError)
			return
		}

		// Stream chunk to final file
		if _, err := io.Copy(finalFile, chunkFile); err != nil {
			chunkFile.Close()
			http.Error(w, "Merge failed", http.StatusInternalServerError)
			return
		}
		chunkFile.Close()
	}

	// 4. Cleanup temp folder
	os.RemoveAll(tempFolderPath)

	w.Write([]byte("File merged successfully"))
}
