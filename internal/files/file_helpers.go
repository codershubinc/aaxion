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

	isAllowed := false
	if strings.HasPrefix(path, getRootPath()) {
		isAllowed = true
	}

	// Allow external storage mount points
	if !isAllowed {
		externalPrefixes := []string{"/media/", "/mnt/", "/run/media/"}
		for _, prefix := range externalPrefixes {
			if strings.HasPrefix(path, prefix) {
				isAllowed = true
				break
			}
		}
	}

	if !isAllowed {
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
