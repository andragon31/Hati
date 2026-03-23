package main

import (
	"fmt"
	"github.com/andragon31/hati/internal/utils"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [agent]",
	Short: "Setup Hati for an AI agent",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agent := args[0]
		log.Info("Setup for Hati not fully implemented yet", "agent", agent)
		fmt.Printf("Generic Hati MCP setup:\nCommand: %s\nArgs: mcp\n", utils.ResolveBinaryPath())
	},
}
