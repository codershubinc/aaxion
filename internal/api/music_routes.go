package api

import (
	"aaxion/internal/music"
	"aaxion/internal/ws"
	"net/http"
)

func AddMusicRoutes() {
	http.HandleFunc("/music/add", music.AddTrackApi)
	http.HandleFunc("/music/search", music.SearchTracksApi)
	http.HandleFunc("/music/all", music.GetTracksApi)
	http.HandleFunc("/music/get", music.GetTrackByIDApi)
	http.HandleFunc("/api/devices", ws.GetDevicesHandler)
	http.HandleFunc("/music/stream", music.StreamTrackApi)
	http.HandleFunc("/music/update", music.UpdateTrackApi)
}
