package routes

import (
	"net/http"
	"time"

	"github.com/PietroPan/music_app/internal/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// SetupRouter sets up all routes and middleware
func SetupRouter(albumHandler *handlers.AlbumHandler, listHandler *handlers.ListHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Albums routes
	r.Get("/albums/search", albumHandler.SearchAlbums)
	r.Get("/albums/{id}", albumHandler.GetAlbum)

	// Lists routes
	r.Post("/lists", listHandler.CreateList)
	r.Get("/lists/{id}", listHandler.GetList)

	// Add album to list
	r.Post("/lists/{id}/albums/{album_id}", listHandler.AddAlbumToList)

	return r
}
