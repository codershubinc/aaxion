package files

import (
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

func viewFiles(dir string) ([]FileInfo, error) {

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
