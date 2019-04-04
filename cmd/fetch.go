package cmd

import (
	"strings"

	"github.com/setplaylistbuilder/builder"
	"github.com/spf13/cobra"
)

var artistName string

var fetchCmd = &cobra.Command{
	Use:   "fetch --artistName <artist name>",
	Short: "Builds a playlist from the latest setlist for the artist",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(artistName) == "" {
			panic("Artist name cannot be empty")
		}
		builder.Build(artistName)
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&artistName, "artist-name", "a", "", "Name of the artist to build playlist based on setlist")

	rootCmd.AddCommand(fetchCmd)
}
