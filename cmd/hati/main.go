package main

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/andragon31/hati/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed plugins/opencode/*
var openCodeFS embed.FS

var (
	version = "0.1.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "hati",
		Short: "Hati - AI Execution & Action Layer",
		Long: `Hati is the execution layer for AI software development.
It complements Skoll by providing the underlying action mechanisms and
state management for complex AI-driven tasks.`,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("data-dir", "d", "", "Data directory (default: ~/.hati)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level: debug|info|warn|error")

	viper.BindPFlag("data_dir", rootCmd.PersistentFlags().Lookup("data-dir"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hati v%s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(tuiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func initConfig() {
	viper.SetDefault("data_dir", utils.GetDefaultDataDir())
	viper.SetDefault("log_level", "info")
}
