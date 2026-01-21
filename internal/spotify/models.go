package spotify

type Album struct {
	SpotifyID   string   `json:"spotify_id"`
	Name        string   `json:"name"`
	ImageURL    string   `json:"image_url"`
	AlbumURL    string   `json:"album_url"`
	ReleaseDate string   `json:"release_date"`
	Artists     []string `json:"artists"`
}

type albumAPIResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ReleaseDate  string `json:"release_date"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
}

type SearchAlbumsResponse struct {
	Albums struct {
		Items []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			ReleaseDate  string `json:"release_date"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"albums"`
}
