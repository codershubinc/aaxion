package db

import "aaxion/internal/models"

func AddTrack(t models.Track) error {
	query := `INSERT INTO tracks (title, artist, album, duration, release_year, file_path, image_path, size) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := GetDB().Exec(query, t.Title, t.Artist, t.Album, t.Duration, t.ReleaseYear, t.FilePath, t.ImagePath, t.Size)
	return err
}

func GetAllTracks() ([]models.Track, error) {
	query := `SELECT id, title, artist, album, duration, release_year, file_path, COALESCE(image_path, '') as image_path, size, created_at FROM tracks`
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		err := rows.Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.Duration, &t.ReleaseYear, &t.FilePath, &t.ImagePath, &t.Size, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func GetTrackByID(id int64) (models.Track, error) {
	query := `SELECT id, title, artist, album, duration, release_year, file_path, COALESCE(image_path, '') as 	image_path, size, created_at FROM tracks WHERE id = ?`
	var t models.Track
	err := GetDB().QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.Duration, &t.ReleaseYear, &t.FilePath, &t.ImagePath, &t.Size, &t.CreatedAt)
	return t, err
}

func SearchTracksByTitle(title string) ([]models.Track, error) {
	query := `SELECT id, title, artist, album, duration, release_year, file_path, COALESCE(image_path, '') as image_path, size, created_at FROM tracks WHERE title LIKE ?`
	rows, err := GetDB().Query(query, "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		err := rows.Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.Duration, &t.ReleaseYear, &t.FilePath, &t.ImagePath, &t.Size, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}
