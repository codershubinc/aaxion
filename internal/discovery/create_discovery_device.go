package discovery

import (
	"aaxion/internal/db"
	"database/sql"
	"log"
)

// var authTokensTableSchema = `
// CREATE TABLE IF NOT EXISTS auth_tokens (
// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
// 	user_id INTEGER NOT NULL,
// 	token TEXT NOT NULL UNIQUE,
// 	type TEXT  DEFAULT 'drive',
// 	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 	FOREIGN KEY(user_id) REFERENCES users(id)
// );
// `

type DeviceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateDiscoveryDevice(name string) (deviceId string, err error) {
	id, err := db.CreateToken()
	if err != nil {
		return "", err
	}
	query := "INSERT INTO discovery_devices (id,name) VALUES (?, ?)"
	_, err = db.GetDB().Exec(query, id, name)
	if err != nil {
		log.Println("Got err", err)
		return "", err
	}

	return id, nil
}

func GetDiscoveryDevices() (dI DeviceInfo, err error) {
	deviceInfo := DeviceInfo{}

	// 1. Use QueryRow for single records
	// 2. Added ORDER BY for consistent results
	// 3. Added LIMIT 1 for performance
	query := "SELECT id, name FROM discovery_devices ORDER BY id ASC LIMIT 1"

	err = db.GetDB().QueryRow(query).Scan(&deviceInfo.ID, &deviceInfo.Name)

	if err != nil {
		// Handle the case where the table is empty
		if err == sql.ErrNoRows {
			// Return empty struct and nil error (or handle as you see fit)
			return DeviceInfo{}, nil
		}
		return DeviceInfo{}, err
	}

	return deviceInfo, nil
}
