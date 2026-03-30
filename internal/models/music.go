package models

import (
	"time"
)

type Track struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Artist      string    `json:"artist" db:"artist"`
	Album       string    `json:"album" db:"album"`
	Duration    float64   `json:"duration" db:"duration"`
	ReleaseYear int       `json:"releaseYear" db:"release_year"`
	FilePath    string    `json:"filePath" db:"file_path"`
	ImagePath   string    `json:"imagePath" db:"image_path"`
	Size        int64     `json:"size" db:"size"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type YTDLPPlaylist struct {
	Type    string `json:"_type"`
	ID      string `json:"id"`
	Entries []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"entries"`
}
