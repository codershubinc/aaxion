//go:build windows

package system

import (
	"aaxion/internal/helpers"
	"net/http"
	"os"
	"path/filepath"
)

func GetSystemStorage(w http.ResponseWriter, r *http.Request) {
	// Dummy response for Windows for now
	response := map[string]interface{}{
		"total":            0,
		"used":             0,
		"available":        0,
		"usage_percentage": 0.0,
		"external_devices": []map[string]interface{}{},
	}

	helpers.SetJSONResponce(w, response)
}

func GetSystemRootPath(w http.ResponseWriter, r *http.Request) {
	rootPath := getRootPath()
	helpers.SetJSONResponce(w, map[string]string{"root_path": rootPath})
}

func getRootPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(string(homeDir))
}
