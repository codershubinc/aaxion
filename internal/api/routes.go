package api

import (
	"aaxion/internal/anonymous_upload"
	"aaxion/internal/auth"
	"aaxion/internal/files"
	img "aaxion/internal/image"
	"aaxion/internal/music"
	"aaxion/internal/music/stats"
	"aaxion/internal/streamer"
	"aaxion/internal/streamer/movies"
	"aaxion/internal/streamer/series"
	sys "aaxion/internal/system"
	"aaxion/internal/tokens"
	"aaxion/internal/webdav"
	"aaxion/internal/ws"
	"net/http"
)

func RegisterRoutes() {
	// ==========================================
	// Authentication
	// ==========================================
	http.HandleFunc("/api/v1/auth/register", auth.Register)
	http.HandleFunc("/api/v1/auth/login", auth.Login)
	http.HandleFunc("/api/v1/auth/logout", auth.Logout)
	// Tokens
	http.HandleFunc("/api/v1/auth/token/generate", auth.AuthMiddleware(tokens.CreateAccessToken))
	http.HandleFunc("/api/v1/auth/token/remove", auth.AuthMiddleware(tokens.CleanAccessTokens))

	// ==========================================
	// WebDAV
	// ==========================================
	webdavHandler := webdav.NewHandler(webdav.GetRootPath())
	http.Handle("/webdav/", webdavHandler)

	// ==========================================
	// Files & Directories
	// ==========================================
	http.HandleFunc("/api/v1/files/view", auth.AuthMiddleware(files.ViewContent))
	http.HandleFunc("/api/v1/files/directory/create", auth.AuthMiddleware(files.CreateDirectory))
	http.HandleFunc("/api/v1/files/unzip", auth.AuthMiddleware(files.UnzipHandler))

	// File Upload (Standard & Chunked)
	http.HandleFunc("/api/v1/files/upload", auth.AuthMiddleware(files.UploadFile))
	http.HandleFunc("/api/v1/files/upload/chunk/start", auth.AuthMiddleware(files.HandleStartChunkUpload))
	http.HandleFunc("/api/v1/files/upload/chunk", auth.AuthMiddleware(files.HandleUploadChunk))
	http.HandleFunc("/api/v1/files/upload/chunk/complete", auth.AuthMiddleware(files.HandleCompleteUpload))

	// File Download & Streaming
	http.HandleFunc("/api/v1/files/download", auth.AuthMiddleware(files.DownloadFileApi))
	http.HandleFunc("/api/v1/files/stream", auth.AuthMiddleware(streamer.StreamFileByPathApi))

	// Images
	http.HandleFunc("/api/v1/images/thumbnail", auth.AuthMiddleware(img.ServeThumbnail))
	http.HandleFunc("/api/v1/images/view", auth.AuthMiddleware(img.ViewImage))

	// ==========================================
	// Temporary File Sharing
	// ==========================================
	http.HandleFunc("/api/v1/share/temp/{token}", files.FileTempShare)
	http.HandleFunc("/api/v1/share/temp/request", files.RequestFileTempShare)

	// ==========================================
	// Anonymous & Token-based Uploads
	// ==========================================
	// Token generation and management (Admin)
	http.HandleFunc("/api/v1/anonymous/token/generate", auth.AuthMiddleware(anonymous_upload.GenerateTokenHandler))
	http.HandleFunc("/api/v1/anonymous/token/revoke", auth.AuthMiddleware(anonymous_upload.RevokeTokenHandler))
	http.HandleFunc("/api/v1/anonymous/tokens", auth.AuthMiddleware(anonymous_upload.ListTokensHandler))
	http.HandleFunc("/api/v1/anonymous/token/info", auth.AuthMiddleware(anonymous_upload.GetTokenInfoHandler))

	// Uploads using token (No Auth Middleware needed, valid token validated in handler)
	http.HandleFunc("/api/v1/anonymous/upload", anonymous_upload.TokenUploadFile)
	http.HandleFunc("/api/v1/anonymous/upload/chunk/start", anonymous_upload.TokenHandleStartChunkUpload)
	http.HandleFunc("/api/v1/anonymous/upload/chunk", anonymous_upload.TokenHandleUploadChunk)
	http.HandleFunc("/api/v1/anonymous/upload/chunk/complete", anonymous_upload.TokenHandleCompleteUpload)
	http.HandleFunc("/api/v1/anonymous/token/validate", anonymous_upload.ValidateTokenHandler)

	// Initialize token cleanup
	anonymous_upload.Initialize()

	// ==========================================
	// System Information
	// ==========================================
	http.HandleFunc("/api/v1/system/info", sys.GetSystemInfoApi) // Unprotected info endpoint
	http.HandleFunc("/api/v1/system/root-path", auth.AuthMiddleware(sys.GetSystemRootPath))
	http.HandleFunc("/api/v1/system/storage", auth.AuthMiddleware(sys.GetSystemStorage))

	// ⚠️⚠️⚠️⚠️ Keep attention : (this endpoints , im not gonna develop further)⚠️⚠️⚠️⚠️
	// ==========================================
	// Movies
	// ==========================================
	http.HandleFunc("/api/v1/movies", auth.AuthMiddleware(movies.ListMoviesApi))
	http.HandleFunc("/api/v1/movies/search", auth.AuthMiddleware(movies.SearchMoviesApi))
	http.HandleFunc("/api/v1/movies/add", auth.AuthMiddleware(movies.AddMovieApi))
	http.HandleFunc("/api/v1/movies/edit", auth.AuthMiddleware(movies.EditMovieApi))
	http.HandleFunc("/api/v1/movies/stream", auth.AuthMiddleware(movies.StreamMovieApi))

	// ==========================================
	// Series & Episodes
	// ==========================================
	http.HandleFunc("/api/v1/series", auth.AuthMiddleware(series.ListSeriesApi))
	http.HandleFunc("/api/v1/series/search", auth.AuthMiddleware(series.SearchSeriesApi))
	http.HandleFunc("/api/v1/series/add", auth.AuthMiddleware(series.AddSeriesApi))
	http.HandleFunc("/api/v1/series/edit", auth.AuthMiddleware(series.EditSeriesApi))

	http.HandleFunc("/api/v1/series/episodes", auth.AuthMiddleware(series.ListEpisodesApi))
	http.HandleFunc("/api/v1/series/episodes/add", auth.AuthMiddleware(series.AddEpisodeApi))
	http.HandleFunc("/api/v1/series/episodes/stream", auth.AuthMiddleware(series.StreamEpisodeApi))

	// ==========================================
	// Music
	// ==========================================
	http.HandleFunc("/api/v1/music", music.GetTracksApi)
	http.HandleFunc("/api/v1/music/get", music.GetTrackByIDApi)
	http.HandleFunc("/api/v1/music/search", music.SearchTracksApi)
	http.HandleFunc("/api/v1/music/add", music.AddTrackApi)
	http.HandleFunc("/api/v1/music/update", music.UpdateTrackApi)
	http.HandleFunc("/api/v1/music/stream", music.StreamTrackApi)

	// Music Stats
	http.HandleFunc("/api/v1/music/stats/play", stats.RecordPlayApi)
	http.HandleFunc("/api/v1/music/stats/favorite", stats.SetFavoriteApi)
	http.HandleFunc("/api/v1/music/stats/play-state", stats.GetPlayStateApi)
	http.HandleFunc("/api/v1/music/stats/play-states", stats.GetAllPlayStatesApi)
	http.HandleFunc("/api/v1/music/stats/is-favorite", stats.CheckFavoriteApi)
	http.HandleFunc("/api/v1/music/stats/last-played", stats.GetLastPlayedApi)

	// ==========================================
	// WebSockets & Devices
	// ==========================================
	http.HandleFunc("/api/v1/devices", ws.GetDevicesHandler)
	http.HandleFunc("/api/v1/ws", ws.Handler)

	// ⚠️⚠️⚠️⚠️ Keep attention : (this endpoints , im not gonna develop further)⚠️⚠️⚠️⚠️
}
