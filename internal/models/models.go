package models

import "time"

type Token struct {
	ID        int64     `json:"id"`
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	FilePath  string    `json:"file_path"`
	Expiry    time.Time `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type Movie struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	FileID      int64     `json:"file_id"`
	CreatedAt   time.Time `json:"created_at"`
	FilePath    string    `json:"file_path"`
	Description string    `json:"description"`
	PosterPath  string    `json:"poster_path"`
	Size        int64     `json:"size"`
	MimeType    string    `json:"mime_type"`
}

type Series struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Episode struct {
	ID            int64     `json:"id"`
	SeriesID      int64     `json:"series_id"`
	FileID        int64     `json:"file_id"`
	SeasonNumber  int       `json:"season_number"`
	EpisodeNumber int       `json:"episode_number"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	FilePath      string    `json:"file_path"`
	Size          int64     `json:"size"`
	MimeType      string    `json:"mime_type"`
	CreatedAt     time.Time `json:"created_at"`
}
