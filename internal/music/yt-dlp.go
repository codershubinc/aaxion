package music

import (
	"aaxion/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DownloadYouTubeAudio(url, musicDir string) (filePath string, err error) {
	args := []string{
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--write-thumbnail",
		"--convert-thumbnails", "png",
		"--embed-thumbnail",
		"--add-metadata",
		"--write-info-json",
		"--print", "after_move:filepath",
		"--progress",
		"--output", filepath.Join(musicDir, "%(title)s.%(ext)s"),
		url,
	}
	cmd := exec.Command("yt-dlp", args...)

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	filePath = strings.TrimSpace(outBuf.String())

	return filePath, nil
}

func ExtractYouTubeMetadata(path string) (track models.Track, err error) {

	baseName := strings.TrimSuffix(path, filepath.Ext(path))
	jsonPath := baseName + ".info.json"
	imagePath := baseName + ".png"

	fileData, err := os.ReadFile(jsonPath)
	if err != nil {
		return models.Track{}, fmt.Errorf("failed to read metadata file: %v", err)
	}

	var info struct {
		Title       string  `json:"title"`
		Artist      string  `json:"artist"`
		Uploader    string  `json:"uploader"`
		Album       string  `json:"album"`
		Duration    float64 `json:"duration"`
		YtUri       string  `json:"webpage_url"`
		ReleaseYear float64 `json:"release_year"`
	}

	if err := json.Unmarshal(fileData, &info); err != nil {
		return models.Track{}, fmt.Errorf("failed to parse metadata: %v", err)
	}

	artist := info.Artist
	if artist == "" {
		artist = info.Uploader
	}
	if artist == "" {
		artist = "Unknown Artist"
	}

	album := info.Album
	if album == "" {
		album = "Unknown Album"
	}

	title := info.Title
	if title == "" {
		title = filepath.Base(baseName)
	}

	var fileSize int64 = 0
	if stat, err := os.Stat(path); err == nil {
		fileSize = stat.Size()
	}
	defer os.Remove(jsonPath)
	return models.Track{
		Title:       title,
		Artist:      artist,
		Album:       album,
		Duration:    info.Duration,
		YtUri:       info.YtUri,
		ReleaseYear: int(info.ReleaseYear),
		FilePath:    path,
		ImagePath:   imagePath,
		Size:        fileSize}, nil
}

func ManageDownloadYoutubePlaylist(url string) ([]string, error) {

	args := []string{
		"--flat-playlist",
		"--dump-single-json",
		url,
	}

	cmd := exec.Command("yt-dlp", args...)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to scout URL: %v", err)
	}

	var data models.YTDLPPlaylist
	if err := json.Unmarshal(outBuf.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("failed to parse scout data: %v", err)
	}

	var urls []string
	if data.Type == "playlist" || data.Type == "multi_video" {
		for _, entry := range data.Entries {
			if entry.ID != "" {
				urls = append(urls, "https://www.youtube.com/watch?v="+entry.ID)
			}
		}
	} else {
		urls = append(urls, url)
	}
	return urls, nil
}

func CleanYouTubeURL(rawURI string) (uri string, typeOfUri string) {
	parsedURL, err := url.Parse(rawURI)
	if err != nil {
		return rawURI, "unknown"
	}

	query := parsedURL.Query()

	if videoID := query.Get("v"); videoID != "" {
		return "https://www.youtube.com/watch?v=" + videoID, "video"
	}

	if playlistID := query.Get("list"); playlistID != "" && strings.Contains(parsedURL.Path, "playlist") {
		return "https://www.youtube.com/playlist?list=" + playlistID, "playlist"
	}

	return rawURI, "unknown"
}
