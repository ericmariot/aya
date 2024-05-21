/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearConfigCmd = &cobra.Command{
	Use:   "clearConfig",
	Short: "Clear the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		err := clearConfig()
		if err != nil {
			fmt.Println("error clearing config file:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clearConfigCmd)
}

func clearConfig() error {
	defaultConfig := Config{
		City:            "",
		LastIP:          "",
		CityCoordinates: make(map[string]CityCoordinates),
	}

	err := saveConfig(defaultConfig)
	if err != nil {
		return fmt.Errorf("error clearing config file: %w", err)
	}
	fmt.Println("Config file cleared.")
	return nil
}
