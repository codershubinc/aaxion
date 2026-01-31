package movies

import (
	"aaxion/internal/db"
	"database/sql"
	"log"
	"time"
)

type Movie struct {
	ID          int64
	Title       string
	FileID      int64
	CreatedAt   time.Time
	FilePath    string
	Description string
	PosterPath  string
	Size        int64
	MimeType    string
}

func AddMovie(title string, fileID int64, filePath, description, posterPath string, size int64, mimeType string) error {
	stmt := `INSERT INTO movies (title, file_id, file_path, description, poster_path, size, mime_type) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.GetDB().Exec(stmt, title, fileID, filePath, description, posterPath, size, mimeType)
	if err != nil {
		log.Printf("Error inserting movie: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Movie added with ID: %d", id)
	return nil
}

func UpdateMovie(id int64, title, description, posterPath string) error {
	stmt := `UPDATE movies SET title = ?, description = ?, poster_path = ? WHERE id = ?`
	_, err := db.GetDB().Exec(stmt, title, description, posterPath, id)
	if err != nil {
		log.Printf("Error updating movie: %v", err)
		return err
	}
	log.Printf("Movie updated with ID: %d", id)
	return nil
}

func GetMovieByID(id int64) (*Movie, error) {
	stmt := `SELECT id, title, file_id, created_at, file_path, description, poster_path, size, mime_type FROM movies WHERE id = ?`
	row := db.GetDB().QueryRow(stmt, id)

	var m Movie
	var createdAt string
	var description, posterPath, mimeType sql.NullString
	var size sql.NullInt64

	err := row.Scan(&m.ID, &m.Title, &m.FileID, &createdAt, &m.FilePath, &description, &posterPath, &size, &mimeType)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		m.Description = description.String
	}
	if posterPath.Valid {
		m.PosterPath = posterPath.String
	}
	if size.Valid {
		m.Size = size.Int64
	}
	if mimeType.Valid {
		m.MimeType = mimeType.String
	}

	m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if m.CreatedAt.IsZero() {
		m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	}

	return &m, nil
}

func ListMovies() ([]Movie, error) {
	stmt := `SELECT id, title, file_id, created_at, file_path, description, poster_path, size, mime_type FROM movies`
	rows, err := db.GetDB().Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		var createdAt string
		var description, posterPath, mimeType sql.NullString
		var size sql.NullInt64

		err := rows.Scan(&m.ID, &m.Title, &m.FileID, &createdAt, &m.FilePath, &description, &posterPath, &size, &mimeType)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			m.Description = description.String
		}
		if posterPath.Valid {
			m.PosterPath = posterPath.String
		}
		if size.Valid {
			m.Size = size.Int64
		}
		if mimeType.Valid {
			m.MimeType = mimeType.String
		}

		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if m.CreatedAt.IsZero() {
			m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		}

		movies = append(movies, m)
	}

	return movies, nil
}

func SearchMovies(query string) ([]Movie, error) {
	stmt := `SELECT id, title, file_id, created_at, file_path, description, poster_path, size, mime_type 
            FROM movies 
            WHERE title LIKE ? OR description LIKE ? 
            ORDER BY created_at DESC`

	searchQuery := "%" + query + "%"
	rows, err := db.GetDB().Query(stmt, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		var createdAt string
		var description, posterPath, mimeType sql.NullString
		var size sql.NullInt64

		if err := rows.Scan(&m.ID, &m.Title, &m.FileID, &createdAt, &m.FilePath, &description, &posterPath, &size, &mimeType); err != nil {
			return nil, err
		}

		if description.Valid {
			m.Description = description.String
		}
		if posterPath.Valid {
			m.PosterPath = posterPath.String
		}
		if size.Valid {
			m.Size = size.Int64
		}
		if mimeType.Valid {
			m.MimeType = mimeType.String
		}

		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if m.CreatedAt.IsZero() {
			m.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		}

		movies = append(movies, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
