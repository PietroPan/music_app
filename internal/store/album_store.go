package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PietroPan/music_app/internal/models"
)

type AlbumStore struct {
	db *sql.DB
}

func NewAlbumStore(db *sql.DB) *AlbumStore {
	return &AlbumStore{db: db}
}

func (s *AlbumStore) GetByID(ctx context.Context, id int64) (*models.Album, error) {
	const query = `
		SELECT id, name, image_url, spotify_id, album_url, rating
		FROM albums
		WHERE id = ?
	`

	var a models.Album
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.Name,
		&a.ImageURL,
		&a.SpotifyID,
		&a.AlbumURL,
		&a.Rating,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("album not found")
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// album_store.go
func (s *AlbumStore) GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Album, error) {
	const query = `
		SELECT id, name, image_url, spotify_id, album_url, rating
		FROM albums
		WHERE spotify_id = ?
	`
	var a models.Album
	err := s.db.QueryRowContext(ctx, query, spotifyID).Scan(
		&a.ID, &a.Name, &a.ImageURL, &a.SpotifyID, &a.AlbumURL, &a.Rating,
	)
	if err == sql.ErrNoRows {
		return nil, nil // not found
	}
	return &a, err
}

func (s *AlbumStore) Create(ctx context.Context, a *models.Album) (*models.Album, error) {
	query := `
		INSERT INTO albums (name, image_url, spotify_id, album_url, rating)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := s.db.QueryRowContext(ctx, query, a.Name, a.ImageURL, a.SpotifyID, a.AlbumURL, a.Rating).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert album: %w", err)
	}

	a.ID = id
	return a, nil
}
