package files

import (
	"aaxion/internal/helpers"
	"net/http"
)

func ViewContent(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		http.Error(w, "Missing 'dir' query parameter", http.StatusBadRequest)
		return
	}
	isSuspicious := expelDotPath(dir)
	if isSuspicious {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}
	if dir == "/" {
		dir = getRootPath()
	}

	files, err := viewContent(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SetJSONResponce(w, files)
}

func CreateDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}
	isSuspicious := expelDotPath(path)
	if isSuspicious {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}
	err := createDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	targetDir := r.URL.Query().Get("dir")
	if targetDir == "" {
		http.Error(w, "Missing 'target_dir' query parameter", http.StatusBadRequest)
		return
	}
	isSuspicious := expelDotPath(targetDir)
	if isSuspicious {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}

	err := uploadLargeFileToDir(w, r, targetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
