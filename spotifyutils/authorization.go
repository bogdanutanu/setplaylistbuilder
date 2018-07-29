package spotifyutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// SpotifyAuthorizedClientBuilder ... Provides a Spotify client with the proper Oauth
// token to access and operate
// a user's account
type SpotifyAuthorizedClientBuilder interface {
	GetSpotifyAuthorizedClient(w http.ResponseWriter) *spotify.Client
}

type spotifyAuthorizedClientBuilderStruct struct {
	auth spotify.Authenticator
	ch   chan *spotify.Client
	// state should be random
	state string
}

func NewSpotifyAuthorizedClientBuilder(redirectURI string) *spotifyAuthorizedClientBuilderStruct {
	return &spotifyAuthorizedClientBuilderStruct{
		auth:  spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate),
		ch:    make(chan *spotify.Client),
		state: "thisshouldberandom",
	}
}

func (acb *spotifyAuthorizedClientBuilderStruct) GetSpotifyAuthorizedClient(w http.ResponseWriter) *spotify.Client {
	var client *spotify.Client
	token, err := acb.getOauthToken()
	if err != nil {
		url := acb.auth.AuthURL(acb.state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	} else {
		acb.buildClient(token, w)
	}
	// wait for auth to complete
	client = <-acb.ch

	return client
}

func (acb *spotifyAuthorizedClientBuilderStruct) getOauthToken() (*oauth2.Token, error) {
	oauthTokenRawJSON, err := ioutil.ReadFile("./oauth-token.json")
	if err != nil {
		return nil, err
	}

	oauth2Token := &oauth2.Token{}
	err = json.Unmarshal(oauthTokenRawJSON, oauth2Token)
	if err != nil {
		return nil, err
	}

	return oauth2Token, nil
}

func (acb *spotifyAuthorizedClientBuilderStruct) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := acb.auth.Token(acb.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != acb.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, acb.state)
	}
	// use the token to get an authenticated client
	// Persist the token to a file in Json format. DO NOT use the channel when the
	// token can be rebuilt from the file.
	acb.buildClient(tok, w)
}

func (acb *spotifyAuthorizedClientBuilderStruct) buildClient(tok *oauth2.Token, w http.ResponseWriter) {
	client := acb.auth.NewClient(tok)
	fmt.Fprintf(w, fmt.Sprintf("Login Completed! Token received: %+v", tok))
	acb.ch <- &client
}
