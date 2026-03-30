package music

import (
	"os"
	"path/filepath"
)

var DIR = filepath.Join(os.Getenv("HOME"), ".aaxion", "music")
