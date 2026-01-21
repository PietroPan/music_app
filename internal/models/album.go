package models

type Album struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	ImageURL  string  `json:"image_url"`
	SpotifyID string  `json:"spotify_id"`
	AlbumURL  string  `json:"album_url"`
	Rating    float64 `json:"rating"`
}

type List struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListWithAlbums struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Albums []Album `json:"albums"`
}
