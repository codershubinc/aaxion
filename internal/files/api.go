package files

import (
	"aaxion/internal/helpers"
	"net/http"
)

func ViewFiles(w http.ResponseWriter, r *http.Request) {
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

	files, err := viewFiles(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SetJSONResponce(w, files)
}
