package image

import (
	"aaxion/internal/files"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
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

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
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

	// Encode and serve as JPEG
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day

	err = jpeg.Encode(w, thumb, &jpeg.Options{Quality: 75})
	if err != nil {
		http.Error(w, "Failed to encode thumbnail", http.StatusInternalServerError)
	}
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
