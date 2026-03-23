package main

import (
	"github.com/andragon31/hati/internal/generator"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Hati in current project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generator.InitProject(".")
	},
}
