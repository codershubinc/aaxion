package api

import (
	"aaxion/internal/anonymous_upload"
	"aaxion/internal/auth"
	"aaxion/internal/files"
	img "aaxion/internal/image"
	"aaxion/internal/streamer/movies"
	"aaxion/internal/streamer/series"
	sys "aaxion/internal/system"
	"aaxion/internal/webdav"
	"net/http"
)

func RegisterRoutes() {

	// Auth routes
	http.HandleFunc("/auth/register", auth.Register)
	http.HandleFunc("/auth/login", auth.Login)
	http.HandleFunc("/auth/logout", auth.Logout)

	// WebDAV Handler
	webdavHandler := webdav.NewHandler(webdav.GetRootPath())
	http.Handle("/webdav/", webdavHandler)

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

	// Token-based anonymous upload routes
	anonymous_upload.RegisterRoutes()

	// Initialize token cleanup
	anonymous_upload.Initialize()

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

	// this is temp route to serve landing page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/landing.html")
	})

	// Login/Register page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/login.html")
	})

	// Web interface for testing temp file share
	http.HandleFunc("/web", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	// Web interface for token-based uploads
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/anonymous/upload.html")
	})

	// Web interface for token management (client-side auth check)
	http.HandleFunc("/admin/tokens", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/anonymous/token-manager.html")
	})

	// Serve static assets for web interface
	http.HandleFunc("/web/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "web/app.js")
	})
	http.HandleFunc("/web/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "web/styles.css")
	})

	// Serve static assets for anonymous upload
	http.HandleFunc("/web/anonymous/token-upload.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "web/anonymous/token-upload.js")
	})
	http.HandleFunc("/web/anonymous/token-manager.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "web/anonymous/token-manager.js")
	})
	http.HandleFunc("/web/anonymous/auth-helper.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "web/anonymous/auth-helper.js")
	})
}
