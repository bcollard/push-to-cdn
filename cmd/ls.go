package cmd

import (
	"context"
	"fmt"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/bcollard/push-to-cdn/internal/gcs"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls [prefix]",
	Aliases: []string{"list"},
	Short:   "List objects in the bucket (optionally filtered by prefix)",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		bucket, err := cfg.RequireBucket()
		if err != nil {
			return err
		}

		prefix := ""
		if len(args) == 1 {
			prefix = args[0]
		}

		ctx := context.Background()
		client, err := gcs.New(ctx, bucket)
		if err != nil {
			return err
		}
		defer client.Close()

		objs, err := client.List(ctx, prefix)
		if err != nil {
			return err
		}
		for _, o := range objs {
			fmt.Printf("%-20s %10d  %s\n", o.Updated, o.Size, o.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
