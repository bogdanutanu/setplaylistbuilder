package main

import (
	"context"
	"fmt"

	"github.com/jm-duarte/setlistfm"
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
}
