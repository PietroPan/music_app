CREATE TABLE IF NOT EXISTS list_albums (
    list_id INTEGER NOT NULL,
    album_id INTEGER NOT NULL,
    PRIMARY KEY (list_id, album_id),
    FOREIGN KEY (list_id) REFERENCES lists(id),
    FOREIGN KEY (album_id) REFERENCES albums(id)
);
