package music

import (
	"aaxion/internal/db"
)

func AddTrack(uri, dir string) error {

	downloadedFilePath, err := DownloadYouTubeAudio(uri, dir)
	if err != nil {
		return err
	}

	track, err := ExtractYouTubeMetadata(downloadedFilePath)
	if err != nil {
		return err
	}

	err = db.AddTrack(track)
	if err != nil {
		return err
	}
	return nil
}
