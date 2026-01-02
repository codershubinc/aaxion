package files

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func sanitizePath(path string) string {
	dir := filepath.Join(string(os.PathSeparator), "aaxion")
	if path == "" {
		return dir
	}

	return path
}
func getRootPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(string(homeDir))
}

func ExpelDotPath(path string) (isSuspicious bool) {

	if !strings.HasPrefix(path, getRootPath()) {
		return true
	}

	if slices.Contains([]string{"..", "."}, path) {
		return true
	}
	if strings.HasPrefix(path, ".") {
		return true
	}
	return false

}
