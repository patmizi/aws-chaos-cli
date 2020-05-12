package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Viper config location
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "go-cli-boileprlate",
	Short: "This is a Cobra/Viper boilerplate",
	Long:  `This is a Cobra/Viper boilerplate program written by patmizi in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd Run with args: %v\n", args)
	},
}

// Execute will execute the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// cobra.OnInitialize(initConfig)
}
