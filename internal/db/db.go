package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
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

	_, err = dbConn.Exec(tokensTableSchema)
	if err != nil {
		return err
	}
	return nil
}
