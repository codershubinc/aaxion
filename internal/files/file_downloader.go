package files

import "os"

func FileDownloader(path string) (p string, error error) {

	// Check if the file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}
	return path, nil
}
