package cmd

import (
	"fmt"
	"log"

	"github.com/bekzourdk/orders/cmd/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"",
		"config to use",
	)

	err := viper.BindPFlag(
		"config",
		rootCmd.PersistentFlags().Lookup("config"),
	)
	if err != nil {
		log.Fatalf("bind config flag: %v", err)
	}

	rootCmd.AddCommand(
		service.Cmd,
	)
}

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use: "orders",

	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.RunE(service.Cmd, args)
		if err != nil {
			return fmt.Errorf("run service cmd: %w", err)
		}

		return nil
	},
}
