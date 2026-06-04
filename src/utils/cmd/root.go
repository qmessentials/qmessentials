package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "utils",
	Short: "QMEssentials Utils CLI",
	Long:  `A collection of utilities for QMEssentials projects.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Requirements: All parameters will be on the command line; no environment variables.
	// We do not call viper.AutomaticEnv() here.
}

func mustBindPFlag(flag string, cmd *cobra.Command) {
	if err := viper.BindPFlag(flag, cmd.Flags().Lookup(flag)); err != nil {
		panic(fmt.Sprintf("failed to bind flag %q: %v", flag, err))
	}
}
