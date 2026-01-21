package spotify

import (
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

func NewConfig() *Config {
	return &Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}
