package models


type TrackPlayState struct {
	TrackId      int    `json:"track_id"`
	UserId       int    `json:"user_id"`
	PlayCount    int    `json:"play_count"`
	LastPlayedAt string `json:"last_played_at"`
}

type FavoriteTrack struct {
	TrackId      int    `json:"track_id"`
	UserId       int    `json:"user_id"`
	CreatedAt    string `json:"created_at"`
}

type LastPlayedTrack struct {
	UserId    int    `json:"user_id"`
	TrackId   int    `json:"track_id"`
	PlayedAt  string `json:"played_at"`
}

