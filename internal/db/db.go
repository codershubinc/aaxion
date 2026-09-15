package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

func InitDb() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appDir := homeDir + "/.aaxion"
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err := os.MkdirAll(appDir, 0700) // Secure permissions
		if err != nil {
			return err
		}
	}

	dbPath := appDir + "/aaxion.db"
	isNewDb := false

	// create DB file if not exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isNewDb = true
		file, err := os.Create(dbPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	//connect to db
	dbConn, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Println("got an err", err)
		return err
	}

	if isNewDb {
		log.Println("Initializing new database tables...")
	}

	schemas := []string{
		tokensTableSchema,
		usersTableSchema,
		authTokensTableSchema,
		moviesTableSchema,
		seriesTableSchema,
		episodesTableSchema,
		discoveryDevices,
		musicTableSchema,
		favoriteTracksTableSchema,
		playStatesTableSchema,
		lastPlayedTableSchema,
		accessTokensTableSchema,
	}

	for _, schema := range schemas {
		_, err := dbConn.Exec(schema)
		if err != nil {
			log.Println("Error creating table: ", err)
			return err
		}
	}
	
	if isNewDb {
		log.Println("DB initialized successfully")
	}
	return nil
}

func GetDB() *sql.DB {
	return dbConn
}
