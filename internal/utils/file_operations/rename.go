package fileOperations

import (
	"os"
)

func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
