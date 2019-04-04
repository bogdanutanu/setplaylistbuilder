package config

// AppConfig ...
type AppConfig struct {
	SetlistFmAPIKey       string `mapstructure:"setlist_fm_api_key"`
	SpotifyOauthTokenFile string `mapstructure:"spotify_oauth_token_file"`
}

// Config ... Holds the config for the app
var Config AppConfig
