package config

// AppConfig ...
type AppConfig struct {
	SetlistFmAPIKey       string `mapstructure:"setlist_fm_api_key" valid:"required"`
	SpotifyOauthTokenFile string `mapstructure:"spotify_oauth_token_file" valid:"-"`
}

// Config ... Holds the config for the app
var Config AppConfig
