package main

import (
	"context"
	"fmt"

	"github.com/jm-duarte/setlistfm"
)

func main() {

	ctx := context.Background()
	client := setlistfm.NewClient("")
	setListQuery := setlistfm.SetlistQuery{
		ArtistName: "kasabian",
	}

	kasabiansetlists, err := client.SearchForSetlists(ctx, setListQuery)
	if err != nil {
		panic(fmt.Sprintf("Error searching for setlists: %s", err))
	}
	fmt.Printf("Response: %+v", kasabiansetlists)
}
