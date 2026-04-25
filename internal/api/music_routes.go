package api

import (
	"aaxion/internal/music"
	"aaxion/internal/music/stats"
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

	// Stats APIs
	http.HandleFunc("/music/stats/play", stats.RecordPlayApi)
	http.HandleFunc("/music/stats/favorite", stats.SetFavoriteApi)
	http.HandleFunc("/music/stats/play_state", stats.GetPlayStateApi)
	http.HandleFunc("/music/stats/play_states", stats.GetAllPlayStatesApi)
	http.HandleFunc("/music/stats/is_favorite", stats.CheckFavoriteApi)
	http.HandleFunc("/music/stats/last_played", stats.GetLastPlayedApi)
}
