package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	config     *Config
	http       *http.Client
	token      string
	tokenExp   time.Time
	tokenMutex sync.Mutex
}

func New(config *Config) *Client {
	return &Client{
		config: config,
		http:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) getToken(ctx context.Context) error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if time.Now().Before(c.tokenExp) && c.token != "" {
		return nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, _ := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://accounts.spotify.com/api/token",
		bytes.NewBufferString(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("spotify token request failed: %s", resp.Status)
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"` // seconds
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return err
	}

	c.token = tr.AccessToken
	c.tokenExp = time.Now().Add(time.Duration(tr.ExpiresIn-30) * time.Second) // 30s buffer
	return nil
}

func (c *Client) doRequest(ctx context.Context, endpoint string, v interface{}) error {
	if err := c.getToken(ctx); err != nil {
		return err
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("spotify request failed: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) SearchAlbums(ctx context.Context, query string) ([]Album, error) {
	endpoint := "https://api.spotify.com/v1/search?type=album&q=" + url.QueryEscape(query) + "&limit=10"
	var resp SearchAlbumsResponse

	if err := c.doRequest(ctx, endpoint, &resp); err != nil {
		return nil, err
	}

	albums := make([]Album, 0, len(resp.Albums.Items))
	for _, item := range resp.Albums.Items {
		imageURL := ""
		if len(item.Images) > 0 {
			imageURL = item.Images[0].URL
		}

		artists := make([]string, len(item.Artists))
		for i, a := range item.Artists {
			artists[i] = a.Name
		}

		albums = append(albums, Album{
			SpotifyID:   item.ID,
			Name:        item.Name,
			ImageURL:    imageURL,
			AlbumURL:    item.ExternalURLs.Spotify,
			ReleaseDate: item.ReleaseDate,
			Artists:     artists,
		})
	}

	return albums, nil
}

func (c *Client) GetAlbumByID(ctx context.Context, spotifyID string) (*Album, error) {
	endpoint := "https://api.spotify.com/v1/albums/" + spotifyID
	var resp albumAPIResponse

	if err := c.doRequest(ctx, endpoint, &resp); err != nil {
		return nil, err
	}

	imageURL := ""
	if len(resp.Images) > 0 {
		imageURL = resp.Images[0].URL
	}

	artists := make([]string, len(resp.Artists))
	for i, a := range resp.Artists {
		artists[i] = a.Name
	}

	return &Album{
		SpotifyID:   resp.ID,
		Name:        resp.Name,
		ImageURL:    imageURL,
		AlbumURL:    resp.ExternalURLs.Spotify,
		ReleaseDate: resp.ReleaseDate,
		Artists:     artists,
	}, nil
}
