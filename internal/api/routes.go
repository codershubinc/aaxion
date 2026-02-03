package api

import (
	"aaxion/internal/auth"
	"aaxion/internal/files"
	img "aaxion/internal/image"
	"aaxion/internal/streamer/movies"
	"aaxion/internal/streamer/series"
	sys "aaxion/internal/system"
	"net/http"
)

func RegisterRoutes() {

	// Auth routes
	http.HandleFunc("/auth/register", auth.Register)
	http.HandleFunc("/auth/login", auth.Login)
	http.HandleFunc("/auth/logout", auth.Logout)

	// file management operations
	http.HandleFunc("/api/files/view", auth.AuthMiddleware(files.ViewContent))
	http.HandleFunc("/files/create-directory", auth.AuthMiddleware(files.CreateDirectory))

	// file upload  operations
	http.HandleFunc("/files/upload", auth.AuthMiddleware(files.UploadFile))

	// file upload  operations - chunked
	http.HandleFunc("/files/upload/chunk/start", auth.AuthMiddleware(files.HandleStartChunkUpload))
	http.HandleFunc("/files/upload/chunk/complete", auth.AuthMiddleware(files.HandleCompleteUpload))
	http.HandleFunc("/files/upload/chunk", auth.AuthMiddleware(files.HandleUploadChunk))

	// file download operations
	http.HandleFunc("/files/download", auth.AuthMiddleware(files.DownloadFileApi))
	http.HandleFunc("/files/thumbnail", auth.AuthMiddleware(img.ServeThumbnail))
	http.HandleFunc("/files/view-image", auth.AuthMiddleware(img.ViewImage))

	// temp files sharing
	http.HandleFunc("/files/d/t/{token}", files.FileTempShare)
	http.HandleFunc("/files/d/r", files.RequestFileTempShare)

	// system info
	http.HandleFunc("/api/system/get-root-path", auth.AuthMiddleware(sys.GetSystemRootPath))
	http.HandleFunc("/api/system/storage", auth.AuthMiddleware(sys.GetSystemStorage))

	// Movies operations
	http.HandleFunc("/api/movies/search", auth.AuthMiddleware(movies.SearchMoviesApi))
	http.HandleFunc("/api/movies/list", auth.AuthMiddleware(movies.ListMoviesApi))
	http.HandleFunc("/api/movies/add", auth.AuthMiddleware(movies.AddMovieApi))
	http.HandleFunc("/api/movies/edit", auth.AuthMiddleware(movies.EditMovieApi))

	// Series operations
	http.HandleFunc("/api/series/list", auth.AuthMiddleware(series.ListSeriesApi))
	http.HandleFunc("/api/series/search", auth.AuthMiddleware(series.SearchSeriesApi))
	http.HandleFunc("/api/series/add", auth.AuthMiddleware(series.AddSeriesApi))
	http.HandleFunc("/api/series/edit", auth.AuthMiddleware(series.EditSeriesApi))

	// Episode operations
	http.HandleFunc("/api/series/episodes/list", auth.AuthMiddleware(series.ListEpisodesApi))
	http.HandleFunc("/api/series/episodes/add", auth.AuthMiddleware(series.AddEpisodeApi))

	// Streamer operations
	http.HandleFunc("/api/stream/movie", auth.AuthMiddleware(movies.StreamMovieApi))
	http.HandleFunc("/api/stream/episode", auth.AuthMiddleware(series.StreamEpisodeApi))

	// this is temp route to serve index.html for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Web interface for testing streaming
	http.HandleFunc("/web", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/stream_test.html")
	})
}
