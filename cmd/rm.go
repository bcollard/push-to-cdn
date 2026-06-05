package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/bcollard/push-to-cdn/internal/gcs"
	"github.com/spf13/cobra"
)

var rmYes bool

var rmCmd = &cobra.Command{
	Use:     "rm <object> [more-objects...]",
	Aliases: []string{"delete", "del"},
	Short:   "Delete one or more objects from the bucket",
	Long: `Delete one or more objects from the bucket. Lists the targets and asks for
confirmation before deleting unless --yes is passed.

The bucket has versioning disabled, so deletes are not recoverable.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		bucket, err := cfg.RequireBucket()
		if err != nil {
			return err
		}

		if !rmYes {
			fmt.Fprintf(os.Stderr, "About to delete %d object(s) from gs://%s:\n", len(args), bucket)
			for _, name := range args {
				fmt.Fprintf(os.Stderr, "  %s\n", name)
			}
			fmt.Fprint(os.Stderr, "Continue? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(os.Stderr, "aborted")
				return nil
			}
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
	rmCmd.Flags().BoolVarP(&rmYes, "yes", "y", false, "skip confirmation prompt")
	rootCmd.AddCommand(rmCmd)
}
