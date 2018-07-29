package spotifyutils

import "github.com/zmb3/spotify"

// SpotifyAuthorizedClient ... Provides a Spotify client with the proper Oauth
// token to access and operate
// a user's account
type SpotifyAuthorizedClient interface {
	GetSpotifyOauthToken()
}

type spotifyAuthorizedClient struct {
	auth spotify.Authenticator
	ch   chan *spotify.Client
	// state should be random
	state string
}
