package api

import (
	"aaxion/internal/music"
	"aaxion/internal/ws"
	"net/http"
)

func AddMusicRoutes() {
	http.HandleFunc("/api/music/add", music.AddTrackApi)
	http.HandleFunc("/api/music/search", music.SearchTracksApi)
	http.HandleFunc("/api/devices", ws.GetDevicesHandler)
}