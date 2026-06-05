package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pushcdn",
	Short: "Push files to a public Google Cloud Storage bucket fronted by a CDN",
}

func Execute(version, commit, date string) {
	rootCmd.Version = version + " (" + commit + ", " + date + ")"
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
