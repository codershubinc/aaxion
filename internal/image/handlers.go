package image

import (
	"aaxion/internal/files"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Ensure formats are registered
var _ = gif.Decode
var _ = jpeg.Decode
var _ = png.Decode

func ServeThumbnail(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	if files.ExpelDotPath(filePath) {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}

	// Check if file exists and get info
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Cache logic
	cacheDir := filepath.Join("uploads_temp", "thumbnails")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		http.Error(w, "Failed to create cache directory", http.StatusInternalServerError)
		return
	}

	// Generate hash based on file path and modification time
	hash := md5.Sum([]byte(filePath + info.ModTime().String()))
	cacheFilename := hex.EncodeToString(hash[:]) + ".jpg"
	cachePath := filepath.Join(cacheDir, cacheFilename)

	// Serve from cache if exists
	if _, err := os.Stat(cachePath); err == nil {
		log.Println("Served from  cache")
		w.Header().Set("Cache-Control", "public, max-age=604800") // Cache for 7 days
		http.ServeFile(w, r, cachePath)
		return
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Failed to decode image: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate new dimensions (max 200px)
	const maxDim = 200
	bounds := img.Bounds()
	wOrig := bounds.Dx()
	hOrig := bounds.Dy()

	var wNew, hNew int
	if wOrig > hOrig {
		wNew = maxDim
		hNew = (hOrig * maxDim) / wOrig
	} else {
		hNew = maxDim
		wNew = (wOrig * maxDim) / hOrig
	}

	// Resize if original is larger than thumbnail
	var thumb image.Image
	if wOrig > maxDim || hOrig > maxDim {
		thumb = resizeNearest(img, wNew, hNew)
	} else {
		thumb = img
	}

	// Save to cache
	outFile, err := os.Create(cachePath)
	if err != nil {
		http.Error(w, "Failed to create cache file", http.StatusInternalServerError)
		return
	}

	err = jpeg.Encode(outFile, thumb, &jpeg.Options{Quality: 75})
	outFile.Close()

	if err != nil {
		http.Error(w, "Failed to encode thumbnail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=604800") // Cache for 7 days
	http.ServeFile(w, r, cachePath)
}

// resizeNearest implements a simple nearest-neighbor resizing
func resizeNearest(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * dx / width
			srcY := y * dy / height
			newImg.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}
	return newImg
}
