package ws

import (
	"encoding/json"
	"net/http"
)

type APIResponseDevice struct {
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
}

func GetDevicesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	clientsMu.Lock()
	activeDevices := make([]APIResponseDevice, 0, len(clients))
	for id, client := range clients {
		activeDevices = append(activeDevices, APIResponseDevice{
			DeviceID:   id,
			DeviceName: client.DeviceName,
		})
	}
	clientsMu.Unlock()

	json.NewEncoder(w).Encode(activeDevices)
}
