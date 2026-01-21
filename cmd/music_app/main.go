package main

import (
	"log"
	"net/http"

	"github.com/PietroPan/music_app/db"
	"github.com/PietroPan/music_app/internal/handlers"
	"github.com/PietroPan/music_app/internal/handlers/routes"
	"github.com/PietroPan/music_app/internal/spotify"
	"github.com/PietroPan/music_app/internal/store"
)

func main() {
	dbConn, err := db.Open("./albums.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	/*if err := db.SeedAlbums(dbConn); err != nil {
		log.Fatal(err)
	}*/

	// Stores
	albumStore := store.NewAlbumStore(dbConn)
	listStore := store.NewListStore(dbConn)

	// Spotify client
	spotifyClient := spotify.New(spotify.NewConfig())

	// Handlers
	albumHandler := handlers.NewAlbumHandler(albumStore, spotifyClient)
	listHandler := handlers.NewListHandler(listStore, albumStore, spotifyClient)

	// Router
	r := routes.SetupRouter(albumHandler, listHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
