package series

import (
	"aaxion/internal/db"
	"aaxion/internal/models"
	"database/sql"
	"log"
	"time"
)

// Series Functions

func AddSeries(title, description string) error {
	stmt := `INSERT INTO series (title, description) VALUES (?, ?)`
	result, err := db.GetDB().Exec(stmt, title, description)
	if err != nil {
		log.Printf("Error inserting series: %v", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Series added with ID: %d", id)
	return nil
}

func UpdateSeries(id int64, title, description string) error {
	stmt := `UPDATE series SET title = ?, description = ? WHERE id = ?`
	_, err := db.GetDB().Exec(stmt, title, description, id)
	if err != nil {
		log.Printf("Error updating series: %v", err)
		return err
	}
	return nil
}

func GetSeriesByID(id int64) (*models.Series, error) {
	stmt := `SELECT id, title, description, created_at FROM series WHERE id = ?`
	row := db.GetDB().QueryRow(stmt, id)

	var s models.Series
	var description sql.NullString
	var createdAt string

	err := row.Scan(&s.ID, &s.Title, &description, &createdAt)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		s.Description = description.String
	}
	s.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if s.CreatedAt.IsZero() {
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	}
	return &s, nil
}

func ListSeries() ([]models.Series, error) {
	stmt := `SELECT id, title, description, created_at FROM series`
	rows, err := db.GetDB().Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []models.Series
	for rows.Next() {
		var s models.Series
		var description sql.NullString
		var createdAt string

		if err := rows.Scan(&s.ID, &s.Title, &description, &createdAt); err != nil {
			return nil, err
		}
		if description.Valid {
			s.Description = description.String
		}
		s.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if s.CreatedAt.IsZero() {
			s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		}
		seriesList = append(seriesList, s)
	}
	return seriesList, nil
}

func SearchSeries(query string) ([]models.Series, error) {
	stmt := `SELECT id, title, description, created_at FROM series WHERE title LIKE ? OR description LIKE ? ORDER BY created_at DESC`
	searchQuery := "%" + query + "%"
	rows, err := db.GetDB().Query(stmt, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []models.Series
	for rows.Next() {
		var s models.Series
		var description sql.NullString
		var createdAt string

		if err := rows.Scan(&s.ID, &s.Title, &description, &createdAt); err != nil {
			return nil, err
		}
		if description.Valid {
			s.Description = description.String
		}
		s.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if s.CreatedAt.IsZero() {
			s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		}
		seriesList = append(seriesList, s)
	}
	return seriesList, nil
}

// Episode Functions

func AddEpisode(seriesID, fileID int64, seasonNum, episodeNum int, title, description, filePath string, size int64, mimeType string) error {
	stmt := `INSERT INTO episodes (series_id, file_id, season_number, episode_number, title, description, file_path, size, mime_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.GetDB().Exec(stmt, seriesID, fileID, seasonNum, episodeNum, title, description, filePath, size, mimeType)
	if err != nil {
		log.Printf("Error inserting episode: %v", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Episode added with ID: %d", id)
	return nil
}

func GetEpisodeByID(id int64) (*models.Episode, error) {
	stmt := `SELECT id, series_id, file_id, season_number, episode_number, title, description, file_path, size, mime_type, created_at FROM episodes WHERE id = ?`
	row := db.GetDB().QueryRow(stmt, id)

	var e models.Episode
	var title, description, filePath, mimeType sql.NullString
	var size sql.NullInt64
	var createdAt string

	err := row.Scan(&e.ID, &e.SeriesID, &e.FileID, &e.SeasonNumber, &e.EpisodeNumber, &title, &description, &filePath, &size, &mimeType, &createdAt)
	if err != nil {
		return nil, err
	}
	if title.Valid {
		e.Title = title.String
	}
	if description.Valid {
		e.Description = description.String
	}
	if filePath.Valid {
		e.FilePath = filePath.String
	}
	if mimeType.Valid {
		e.MimeType = mimeType.String
	}
	if size.Valid {
		e.Size = size.Int64
	}
	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if e.CreatedAt.IsZero() {
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	}
	return &e, nil
}

func ListEpisodes(seriesID int64) ([]models.Episode, error) {
	stmt := `SELECT id, series_id, file_id, season_number, episode_number, title, description, file_path, size, mime_type, created_at FROM episodes WHERE series_id = ? ORDER BY season_number, episode_number`
	rows, err := db.GetDB().Query(stmt, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var e models.Episode
		var title, description, filePath, mimeType sql.NullString
		var size sql.NullInt64
		var createdAt string

		if err := rows.Scan(&e.ID, &e.SeriesID, &e.FileID, &e.SeasonNumber, &e.EpisodeNumber, &title, &description, &filePath, &size, &mimeType, &createdAt); err != nil {
			return nil, err
		}
		if title.Valid {
			e.Title = title.String
		}
		if description.Valid {
			e.Description = description.String
		}
		if filePath.Valid {
			e.FilePath = filePath.String
		}
		if mimeType.Valid {
			e.MimeType = mimeType.String
		}
		if size.Valid {
			e.Size = size.Int64
		}
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if e.CreatedAt.IsZero() {
			e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		}
		episodes = append(episodes, e)
	}
	return episodes, nil
}
