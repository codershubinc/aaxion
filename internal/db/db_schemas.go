package db

var tokensTableSchema = `
CREATE TABLE IF NOT EXISTS tokens (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	token TEXT NOT NULL UNIQUE,
	token_type TEXT NOT NULL,
	file_path TEXT NOT NULL,
	expiry DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

var usersTableSchema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

var authTokensTableSchema = `
CREATE TABLE IF NOT EXISTS auth_tokens (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	token TEXT NOT NULL UNIQUE,
	type TEXT  DEFAULT 'drive',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id)
);
`
var moviesTableSchema = `
CREATE TABLE IF NOT EXISTS movies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	file_id INTEGER NOT NULL, -- Link to physical file
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	file_path TEXT NOT NULL,
	description TEXT,
	poster_path TEXT,
	size INTEGER,
	mime_type TEXT,
	UNIQUE(title), -- Prevent duplicate movie titles
	FOREIGN KEY(file_id) REFERENCES files(id) ON DELETE CASCADE
);
`

var seriesTableSchema = `
CREATE TABLE IF NOT EXISTS series (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL UNIQUE, -- Unique constraint prevents duplicates
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

var episodesTableSchema = `
CREATE TABLE IF NOT EXISTS episodes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	series_id INTEGER NOT NULL,
	file_id INTEGER NOT NULL,   -- Link to physical file
	season_number INTEGER NOT NULL,
	episode_number INTEGER NOT NULL,
	title TEXT,                 -- Optional: "Pilot", "Ozymandias"
	description TEXT,
	file_path TEXT NOT NULL,
	size INTEGER,
	mime_type TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE,
	FOREIGN KEY(file_id) REFERENCES files(id) ON DELETE CASCADE,
	UNIQUE(series_id, season_number, episode_number) -- Prevent duplicate episodes
);
`
