package cmd

import (
	"fmt"

	"github.com/bcollard/push-to-cdn/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage pushcdn configuration",
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
			fmt.Printf("%-12s %s\n", key+":", val)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available configuration keys",
	Run: func(cmd *cobra.Command, args []string) {
		for _, key := range config.Keys {
			fmt.Println(key)
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration key",
	Args:  cobra.ExactArgs(2),
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

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.Set(key, value)
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("Set %s\n", key)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd, configListCmd, configSetCmd)
	rootCmd.AddCommand(configCmd)
}
