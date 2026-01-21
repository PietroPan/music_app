package db

import "database/sql"

func SeedAlbums(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		image_url TEXT,
		spotify_id TEXT,
		album_url TEXT,
		rating REAL
	);`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	const seed = `
	INSERT INTO albums (name, image_url, spotify_id, album_url, rating)
	SELECT
		'In Rainbows',
		'https://example.com/in-rainbows.jpg',
		'spotify123',
		'https://open.spotify.com/album/xyz',
		4.8
	WHERE NOT EXISTS (SELECT 1 FROM albums);
	`

	_, err := db.Exec(seed)
	return err
}
