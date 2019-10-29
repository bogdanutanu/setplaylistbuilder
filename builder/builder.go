package builder

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jm-duarte/setlistfm"
	"github.com/setplaylistbuilder/config"
	"github.com/setplaylistbuilder/setlist"
	"github.com/setplaylistbuilder/spotifyutils"
	"github.com/zmb3/spotify"
)

func concatSimpleArtistsNames(simpleArtists []spotify.SimpleArtist) string {
	artistNames := make([]string, len(simpleArtists))
	for i, simpleArtist := range simpleArtists {
		artistNames[i] = simpleArtist.Name
	}
	return strings.Join(artistNames, ", ")
}

func extractTrackIDs(fullTracks []spotify.FullTrack) []spotify.ID {
	trackIDs := make([]spotify.ID, len(fullTracks))
	for i, fullTrack := range fullTracks {
		trackIDs[i] = fullTrack.ID
	}
	return trackIDs
}

func Build(artistName string) {
	ctx := context.Background()
	client := setlistfm.NewClient(config.Config.SetlistFmAPIKey)
	setListQuery := setlistfm.SetlistQuery{
		ArtistName: artistName,
	}

	setlists, err := client.SearchForSetlists(ctx, setListQuery)
	if err != nil {
		panic(fmt.Sprintf("Error searching for setlists: %s", err))
	}

	// We will try to find the most recent non-empty setlist from the first
	// results page only
	lastSetlist := setlist.ExtractMostRecent(setlists.Setlists)

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

	fullTracks := make([]spotify.FullTrack, 0)

	for i, set := range lastSetlist.Sets.Set {
		fmt.Printf("Set %d\n", i)
		for j, song := range set.Song {
			fmt.Printf("%2d: %s\n", j, song.Name)

			results, err := spotifyClient.Search(fmt.Sprintf("%s artist:%s", song.Name, artistName), spotify.SearchTypeTrack)
			if err != nil {
				log.Fatalf("Error searching for '%s': %v", song.Name, err)
				continue
			}
			if results.Tracks != nil {
				for i, track := range results.Tracks.Tracks {
					fmt.Printf("\t%s - %s - %v\n", track.Name, track.Album.Name, concatSimpleArtistsNames(track.Artists))
					if i == 0 {
						fullTracks = append(fullTracks, track)
					}
				}
			}
			fmt.Println()
		}
	}

	fmt.Printf("\n Proposed setplaylist:\n\n")

	for i, track := range fullTracks {
		fmt.Printf("\t%2d. %s - %s - %v\n", i+1, track.Name, track.Album.Name, concatSimpleArtistsNames(track.Artists))
	}

	reader := bufio.NewReader(os.Stdin)
	answer := ""
	for answer != "y" && answer != "n" {
		fmt.Print("Proceed with this setlist? (y/n): ")
		answer, _ = reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")
		fmt.Printf("Answer: %s\n", answer)
	}
	switch answer {
	case "y":
	case "n":
		os.Exit(0)
	default:
		log.Fatal("Answer not supported")
		os.Exit(-1)
	}

	fullPlaylist, err := spotifyClient.CreatePlaylistForUser(
		user.ID,
		"Setlist "+artistName,
		"Playlist created by SetPLayListBuilder",
		false)
	if err != nil {
		log.Fatalf("Error creating playlist for user '%s': %v", user.ID, err)
		os.Exit(-1)
	}

	spotifyClient.AddTracksToPlaylist(fullPlaylist.ID, extractTrackIDs(fullTracks)...)

	fmt.Printf("Playlist '%s' wiht ID '%s' created successfully.\n", fullPlaylist.Name, fullPlaylist.ID)
}
