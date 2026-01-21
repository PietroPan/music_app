package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PietroPan/music_app/internal/models"
)

type ListStore struct {
	db *sql.DB
}

func NewListStore(db *sql.DB) *ListStore {
	return &ListStore{db: db}
}

// Create a new list
func (s *ListStore) Create(ctx context.Context, name string) (*models.List, error) {
	const query = `
		INSERT INTO lists (name)
		VALUES (?)
		RETURNING id, name
	`

	var list models.List
	err := s.db.QueryRowContext(ctx, query, name).Scan(&list.ID, &list.Name)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

// Get a list by ID
func (s *ListStore) GetByID(ctx context.Context, id int64) (*models.List, error) {
	const query = `
		SELECT id, name
		FROM lists
		WHERE id = ?
	`

	var list models.List
	err := s.db.QueryRowContext(ctx, query, id).Scan(&list.ID, &list.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("list not found")
	}
	if err != nil {
		return nil, err
	}

	return &list, nil
}

// list_store.go
func (s *ListStore) AddAlbum(ctx context.Context, listID, albumID int64) error {
	const query = `
		INSERT OR IGNORE INTO list_albums (list_id, album_id)
		VALUES (?, ?)
	`
	_, err := s.db.ExecContext(ctx, query, listID, albumID)
	return err
}

// store/list_store.go
func (s *ListStore) GetWithAlbums(ctx context.Context, listID int64) (*models.ListWithAlbums, error) {
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return nil, err
	}

	albums, err := s.GetAlbumsForList(ctx, listID) // new method in store
	if err != nil {
		return nil, err
	}

	var albumInfos []models.Album
	for _, a := range albums {
		albumInfos = append(albumInfos, models.Album{
			ID:        a.ID,
			Name:      a.Name,
			ImageURL:  a.ImageURL,
			SpotifyID: a.SpotifyID,
			AlbumURL:  a.AlbumURL,
			Rating:    a.Rating,
		})
	}

	return &models.ListWithAlbums{
		ID:     list.ID,
		Name:   list.Name,
		Albums: albumInfos,
	}, nil
}

func (s *ListStore) GetAlbumsForList(ctx context.Context, listID int64) ([]models.Album, error) {
	query := `
		SELECT a.id, a.name, a.image_url, a.spotify_id, a.album_url, a.rating
		FROM albums a
		INNER JOIN list_albums la ON la.album_id = a.id
		WHERE la.list_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("query albums for list: %w", err)
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Name, &a.ImageURL, &a.SpotifyID, &a.AlbumURL, &a.Rating); err != nil {
			return nil, fmt.Errorf("scan album: %w", err)
		}
		albums = append(albums, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return albums, nil
}
