package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PietroPan/music_app/internal/spotify"
	"github.com/PietroPan/music_app/internal/store"
	"github.com/go-chi/chi/v5"
)

type AlbumHandler struct {
	store         *store.AlbumStore
	spotifyClient *spotify.Client
}

func NewAlbumHandler(store *store.AlbumStore, spotifyClient *spotify.Client) *AlbumHandler {
	return &AlbumHandler{
		store:         store,
		spotifyClient: spotifyClient,
	}
}

func (h *AlbumHandler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid album id", http.StatusBadRequest)
		return
	}

	album, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func (h *AlbumHandler) SearchAlbums(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "`q` query parameter is required", http.StatusBadRequest)
		return
	}

	albums, err := h.spotifyClient.SearchAlbums(r.Context(), query)
	if err != nil {
		http.Error(w, "failed to search albums: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}
