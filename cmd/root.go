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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd Run with args: %v\n", args)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
	},
}

// Execute will execute the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// cobra.OnInitialize(initConfig)
}
