package cmd

import (
	"fmt"
	"strings"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url <object>",
	Short: "Print the public URL for an object",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		if cfg.BaseURL == "" {
			return fmt.Errorf("no base-url configured — run: pushcdn config set base-url https://cdn.runlocal.dev")
		}
		name := strings.TrimLeft(args[0], "/")
		fmt.Printf("%s/%s\n", strings.TrimRight(cfg.BaseURL, "/"), name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)
}
