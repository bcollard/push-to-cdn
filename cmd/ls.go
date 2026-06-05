package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/bcollard/push-to-cdn/internal/gcs"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls [pattern]",
	Aliases: []string{"list"},
	Short:   "List objects in the bucket, optionally filtered by prefix or glob",
	Long: `List objects in the bucket.

With no argument, lists everything. With a plain string argument, treats it as a
prefix (server-side filter). When the argument contains glob characters ('*',
'?', or '['), it is sent as a server-side glob match:

  *      matches any character except '/'
  **     matches any character including '/'
  ?      matches exactly one character except '/'
  [a-z]  character class

Quote the pattern to keep your shell from expanding it before pushcdn sees it.`,
	Example: `  pushcdn ls                       # everything
  pushcdn ls brand/                # objects under the "brand/" prefix
  pushcdn ls '*gorilla*'           # any object with "gorilla" in its name (no '/')
  pushcdn ls '**gorilla**'         # same, but match across folder boundaries
  pushcdn ls 'logo.???'            # logo.png, logo.svg, …`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		bucket, err := cfg.RequireBucket()
		if err != nil {
			return err
		}

		var q gcs.ListQuery
		if len(args) == 1 {
			arg := args[0]
			if isGlob(arg) {
				q.MatchGlob = arg
			} else {
				q.Prefix = arg
			}
		}

		ctx := context.Background()
		client, err := gcs.New(ctx, bucket)
		if err != nil {
			return err
		}
		defer client.Close()

		objs, err := client.List(ctx, q)
		if err != nil {
			return err
		}
		for _, o := range objs {
			fmt.Printf("%-20s %10d  %s\n", o.Updated, o.Size, o.Name)
		}
		return nil
	},
}

func isGlob(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
