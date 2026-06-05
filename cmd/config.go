package cmd

import (
	"fmt"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage pushcdn configuration",
	Long: `Manage the persistent configuration at ~/.config/pushcdn/config.json.

Three keys are supported:

  bucket     ` + config.Descriptions["bucket"] + `
             example: ` + config.Examples["bucket"] + `

  project    ` + config.Descriptions["project"] + `
             example: ` + config.Examples["project"] + `

  base-url   ` + config.Descriptions["base-url"] + `
             example: ` + config.Examples["base-url"] + `

Resolution order at runtime: env var (PUSHCDN_BUCKET, PUSHCDN_PROJECT,
PUSHCDN_BASE_URL) → stored config → default (base-url is derived from
bucket if unset).`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		for _, key := range config.Keys {
			val := cfg.Get(key)
			if val == "" {
				val = "(not set)"
			}
			fmt.Printf("%-10s %s\n", key+":", val)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available configuration keys with descriptions and examples",
	Run: func(cmd *cobra.Command, args []string) {
		for _, key := range config.Keys {
			fmt.Printf("%s\n", key)
			fmt.Printf("    %s\n", config.Descriptions[key])
			fmt.Printf("    example: %s\n\n", config.Examples[key])
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration key",
	Long: `Set a configuration key. Values are normalized (gs:// prefix stripped from
bucket, trailing / stripped from base-url) and validated before being saved.

Examples:
  pushcdn config set bucket   cdn.runlocal.dev
  pushcdn config set project  my-gcp-project
  pushcdn config set base-url https://cdn.runlocal.dev`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		valid := false
		for _, k := range config.Keys {
			if k == key {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unknown config key %q — run 'pushcdn config list' to see available keys", key)
		}

		normalized, note := config.Normalize(key, value)
		if note != "" {
			fmt.Fprintf(cmd.ErrOrStderr(), "note: %s\n", note)
		}
		if err := config.Validate(key, normalized); err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.Set(key, normalized)
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("Set %s = %s\n", key, normalized)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd, configListCmd, configSetCmd)
	rootCmd.AddCommand(configCmd)
}
