package files

import (
	fileOperations "aaxion/internal/utils/file_operations"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

type UnzipRequest struct {
	ZipPath string `json:"zip_path"`
	DestDir string `json:"dest_dir"`
}

func UnzipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UnzipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	absZipPath := req.ZipPath

	absDestDir := req.DestDir
	if absDestDir == "" {
		absDestDir = strings.TrimSuffix(absZipPath, filepath.Ext(absZipPath))
	}

	fmt.Println("Dest dir :", absDestDir)
	fmt.Println("Zip path :", absZipPath)
	err := fileOperations.ExtractZip(absZipPath, absDestDir)
	if err != nil {
		http.Error(w, "Failed to extract zip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Archive extracted successfully",
		"dest":    absDestDir,
	})
}
