package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

func InitDb() error {
	// create DB file if not exists

	dbPath := ".aaxion.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	//connect to db
	var err error
	dbConn, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Println("got an err", err)
		return err
	}

	//create tables if not exist
	log.Println("Creating tables")

	schemas := []string{
		tokensTableSchema,
		usersTableSchema,
		authTokensTableSchema,
		moviesTableSchema,
		seriesTableSchema,
		episodesTableSchema,
	}

	for _, schema := range schemas {
		_, err := dbConn.Exec(schema)
		if err != nil {
			log.Println("Error creating table: ", err)
			return err
		}
	}
	log.Println("DB initialized successfully")
	return nil
}

func GetDB() *sql.DB {
	return dbConn
}
