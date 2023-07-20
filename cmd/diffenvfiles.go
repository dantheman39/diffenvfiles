package cmd

import (
	"fmt"
	"github.com/dantheman39/diffenvfiles/pkg"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	var ignoreQuotes bool
	command := &cobra.Command{
		Use:          "diffenvfiles envfile1 envfile2",
		Short:        "Compare two dot-env files",
		Long:         "Run this command to show what variables are different between two .env files",
		Version:      "v1.2.0",
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			envPath1 := args[0]
			envPath2 := args[1]
			env1, err1 := os.ReadFile(envPath1)
			env2, err2 := os.ReadFile(envPath2)

			if err1 != nil {
				return fmt.Errorf("there was an error reading envfile1: %w", err1)
			}
			if err2 != nil {
				return fmt.Errorf("there was an error reading envfile2: %w", err2)
			}
			config := pkg.Config{IgnoreQuotes: ignoreQuotes}
			return pkg.DiffEnvFiles(
				pkg.EnvFile{
					Path: envPath1,
					Data: env1,
				},
				pkg.EnvFile{
					Path: envPath2,
					Data: env2,
				},
				config,
			)
		},
	}
	command.Flags().BoolVar(
		&ignoreQuotes,
		"ignore-quotes",
		false,
		"if the only difference between a value is if it's surrounded in 'single quotes' or \"double quotes\", don't consider them different",
	)

	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Println("")
}
