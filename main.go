package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jm-duarte/setlistfm"
	"github.com/setplaylistbuilder/spotifyutils"
	"github.com/spf13/viper"
)

type AppConfig struct {
	SetlistFmAPIKey       string
	SpotifyOauthTokenFile string
}

var cfg AppConfig

func readConfig() {
	viper.SetConfigName("setplaylistbuilder")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("Fatal error parsing config file: %s \n", err))
	}
}

func init() {
	readConfig()
}

func main() {
	ctx := context.Background()
	client := setlistfm.NewClient(cfg.SetlistFmAPIKey)
	setListQuery := setlistfm.SetlistQuery{
		ArtistName: "kasabian",
	}

	kasabiansetlists, err := client.SearchForSetlists(ctx, setListQuery)
	if err != nil {
		panic(fmt.Sprintf("Error searching for setlists: %s", err))
	}
	fmt.Printf("Response: %+v", kasabiansetlists)

	fmt.Println()
	fmt.Println()

	spotifyClientBuilder := spotifyutils.NewSpotifyAuthorizedClientBuilder("http://localhost:8080/callback")

	http.HandleFunc("/callback", spotifyClientBuilder.CompleteAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	spotifyClient := spotifyClientBuilder.GetSpotifyAuthorizedClient()
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	playlists, err := spotifyClient.GetPlaylistsForUser(user.ID)
	if err != nil {
		panic(fmt.Sprintf("Error retrieving user's playlists: %v", err))
	}

	for _, p := range playlists.Playlists {
		fmt.Printf("%s\n", p.Name)
	}

}
