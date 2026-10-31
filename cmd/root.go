package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh tl",
	Short: "Tupynambá Lucas Enterprise Developer CLI Suite",
	Long: `gh tl provides unified developer tooling, automated workstation onboarding,
cryptographic verification, multi-repo synchronization, and organizational utilities
for the tupynambalucas-dev ecosystem.`,
}

// Execute runs the root CLI command suite.
func Execute() error {
	return rootCmd.Execute()
}
