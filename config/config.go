package config

// AppConfig ...
type AppConfig struct {
	SetlistFmAPIKey       string `json:"setlist_fm_api_key"`
	SpotifyOauthTokenFile string `json:"spotify_oauth_token_file"`
}

// Config ... Holds the config for the app
var Config AppConfig
