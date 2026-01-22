CREATE TABLE IF NOT EXISTS albums (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    spotify_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    image_url TEXT,
    album_url TEXT,
    rating REAL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
