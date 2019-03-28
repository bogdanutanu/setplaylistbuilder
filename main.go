package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zmb3/spotify"

	"github.com/jm-duarte/setlistfm"
	"github.com/setplaylistbuilder/setlist"
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
	// fmt.Printf("Kasabian setlists: %+v", kasabiansetlists)
	fmt.Println()
	fmt.Println()

	// We will try to find the most recent non-empty setlist from the first
	// results page only
	lastSetlist := setlist.ExtractMostRecent(kasabiansetlists.Setlists)
	fmt.Printf("Most recent setlist: %+v", lastSetlist)

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

	// playlists, err := spotifyClient.GetPlaylistsForUser(user.ID)
	// if err != nil {
	// 	panic(fmt.Sprintf("Error retrieving user's playlists: %v", err))
	// }
	// for _, p := range playlists.Playlists {
	// 	fmt.Printf("%s\n", p.Name)
	// }

	for i, set := range lastSetlist.Sets.Set {
		fmt.Printf("Set %d\n", i)
		for j, song := range set.Song {
			fmt.Printf("%2d: %s\n", j, song.Name)

			results, err := spotifyClient.Search(fmt.Sprintf("%s artist:kasabian", song.Name), spotify.SearchTypeTrack)
			if err != nil {
				log.Fatalf("Error searching for '%s': %v", song.Name, err)
				continue
			}
			if results.Tracks != nil {
				for _, track := range results.Tracks.Tracks {
					fmt.Printf("\t%s - %s - %v\n", track.Name, track.Album.Name, concatSimpleArtistsNames(track.Artists))
				}
			}
			fmt.Println()
		}
	}

}

func concatSimpleArtistsNames(simpleArtists []spotify.SimpleArtist) string {
	artistNames := make([]string, len(simpleArtists))
	for i, simpleArtist := range simpleArtists {
		artistNames[i] = simpleArtist.Name
	}
	return strings.Join(artistNames, ", ")
}
