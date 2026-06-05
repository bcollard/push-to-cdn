package cmd

import (
	"context"
	"fmt"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/bcollard/push-to-cdn/internal/gcs"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm <object> [more-objects...]",
	Aliases: []string{"delete", "del"},
	Short:   "Delete one or more objects from the bucket",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		bucket, err := cfg.RequireBucket()
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := gcs.New(ctx, bucket)
		if err != nil {
			return err
		}
		defer client.Close()

		for _, name := range args {
			if err := client.Delete(ctx, name); err != nil {
				return fmt.Errorf("deleting %s: %w", name, err)
			}
			fmt.Printf("deleted %s/%s\n", bucket, name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
