package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PietroPan/music_app/internal/models"
	"github.com/PietroPan/music_app/internal/spotify"
	"github.com/PietroPan/music_app/internal/store"
	"github.com/go-chi/chi/v5"
)

// ListHandler handles list-related endpoints
type ListHandler struct {
	store         *store.ListStore
	albumStore    *store.AlbumStore
	spotifyClient *spotify.Client
}

// NewListHandler creates a new ListHandler
func NewListHandler(listStore *store.ListStore, albumStore *store.AlbumStore, spotifyClient *spotify.Client) *ListHandler {
	return &ListHandler{
		store:         listStore,
		albumStore:    albumStore,
		spotifyClient: spotifyClient,
	}
}

// POST /lists
func (h *ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	list, err := h.store.Create(r.Context(), payload.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(list)
}

// GET /lists/:id
func (h *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	listIDStr := chi.URLParam(r, "id")
	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid list id", http.StatusBadRequest)
		return
	}

	listWithAlbums, err := h.store.GetWithAlbums(r.Context(), listID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listWithAlbums)
}

// POST /lists/:id/albums/:album_id
func (h *ListHandler) AddAlbumToList(w http.ResponseWriter, r *http.Request) {
	listIDStr := chi.URLParam(r, "id")
	spotifyID := chi.URLParam(r, "album_id")

	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid list id", http.StatusBadRequest)
		return
	}

	// Check if list exists
	list, err := h.store.GetByID(r.Context(), listID)
	if err != nil {
		http.Error(w, "list not found", http.StatusNotFound)
		return
	}

	// Check if album exists in DB
	album, err := h.albumStore.GetBySpotifyID(r.Context(), spotifyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if album == nil {
		// Album not in DB, fetch from Spotify
		spotifyAlbum, err := h.spotifyClient.GetAlbumByID(r.Context(), spotifyID)
		if err != nil {
			http.Error(w, "failed to fetch album from Spotify: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Map Spotify album to DB album
		albumData := &models.Album{
			Name:      spotifyAlbum.Name,
			ImageURL:  spotifyAlbum.ImageURL,
			SpotifyID: spotifyAlbum.SpotifyID,
			AlbumURL:  spotifyAlbum.AlbumURL,
			Rating:    0.0, // default rating
		}

		// Save to DB
		album, err = h.albumStore.Create(r.Context(), albumData)
		if err != nil {
			http.Error(w, "failed to save album: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Add album to list
	if err := h.store.AddAlbum(r.Context(), list.ID, album.ID); err != nil {
		http.Error(w, "failed to add album to list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"list_id":    list.ID,
		"album_id":   album.ID,
		"spotify_id": album.SpotifyID,
		"message":    "album added to list",
	})
}
