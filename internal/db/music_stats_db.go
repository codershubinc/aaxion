package db

import (
	"aaxion/internal/models"
	"database/sql"
	"log"
)

// RecordPlay updates the play count and last played time for a user and track,
// and optionally sets the user's global last played track.
func RecordPlay(trackId int, userId int, playedAt string) {
	// Update play_states
	playQuery := `
		INSERT INTO play_states (user_id, track_id, play_count, last_played_at) 
		VALUES (?, ?, 1, ?) 
		ON CONFLICT(user_id, track_id) DO UPDATE SET 
			play_count = play_count + 1,
			last_played_at = ?
	`
	_, err := GetDB().Exec(playQuery, userId, trackId, playedAt, playedAt)
	if err != nil {
		log.Println("Got err updating play state", err)
	}

	// Update last_played_tracks (global)
	lastPlayedQuery := `
		INSERT INTO last_played_tracks (user_id, track_id, played_at) 
		VALUES (?, ?, ?) 
		ON CONFLICT(user_id) DO UPDATE SET 
			track_id = ?,
			played_at = ?
	`
	_, err = GetDB().Exec(lastPlayedQuery, userId, trackId, playedAt, trackId, playedAt)
	if err != nil {
		log.Println("Got err updating last played track", err)
	}
}

// SetFavorite toggles or explicitly sets the favorite status.
func SetFavorite(trackId int, userId int, isFavorite bool) error {
	if isFavorite {
		query := "INSERT OR IGNORE INTO favorite_tracks (track_id) VALUES ( ?)"
		_, err := GetDB().Exec(query, trackId)
		if err != nil {
			log.Println("Got err adding favorite", err)
			return err
		}
	} else {
		query := "DELETE FROM favorite_tracks WHERE  track_id = ?"
		_, err := GetDB().Exec(query, userId, trackId)
		if err != nil {
			log.Println("Got err removing favorite", err)
		}

	}
	return nil
}

// GetPlayState returns the play state for a specific user and track.
func GetPlayState(trackId int, userId int) (state models.TrackPlayState, err error) {
	query := "SELECT user_id, track_id, play_count, last_played_at FROM play_states WHERE user_id = ? AND track_id = ?"

	row := GetDB().QueryRow(query, userId, trackId)

	err = row.Scan(&state.UserId, &state.TrackId, &state.PlayCount, &state.LastPlayedAt)
	if err == sql.ErrNoRows {
		state = models.TrackPlayState{UserId: userId, TrackId: trackId, PlayCount: 0}
		err = nil
	}
	return
}

// GetAllPlayStates returns all play states for a specific user.
func GetAllPlayStates(userId int) ([]models.TrackPlayState, error) {
	query := "SELECT user_id, track_id, play_count, last_played_at FROM play_states WHERE user_id = ?"

	rows, err := GetDB().Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []models.TrackPlayState
	for rows.Next() {
		var state models.TrackPlayState
		if err := rows.Scan(&state.UserId, &state.TrackId, &state.PlayCount, &state.LastPlayedAt); err == nil {
			states = append(states, state)
		}
	}

	// Return empty slice instead of null if no states exist
	if states == nil {
		states = []models.TrackPlayState{}
	}

	return states, nil
}

// CheckFavorite checks if a track is favorited by the user.
func CheckFavorite(trackId int, userId int) bool {
	var id int
	query := "SELECT id FROM favorite_tracks WHERE user_id = ? AND track_id = ?"
	err := GetDB().QueryRow(query, userId, trackId).Scan(&id)
	return err == nil
}

// GetUserLastPlayedTrack returns the single most recently played track for a user.
func GetUserLastPlayedTrack(userId int) (state models.LastPlayedTrack, err error) {
	query := "SELECT user_id, track_id, played_at FROM last_played_tracks WHERE user_id = ?"

	row := GetDB().QueryRow(query, userId)

	err = row.Scan(&state.UserId, &state.TrackId, &state.PlayedAt)
	return
}

func DeletePlayState(trackId int, userId int) {
	query := "DELETE FROM play_states WHERE track_id = ? AND user_id = ?"
	_, err := GetDB().Exec(query, trackId, userId)
	if err != nil {
		log.Println("Got delete err", err)
	}
}
