package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tupynambalucas-dev/gh-tl/cmd/onboard"
)

var rootCmd = &cobra.Command{
	Use:   "gh tl",
	Short: "Tupynambá Lucas Enterprise Developer CLI Suite",
	Long: `gh tl is the centralized developer toolchain for the tupynambalucas-dev ecosystem,
orchestrating workstation provisioning, cryptographic compliance, and multi-repo workflows.`,
}

// Execute runs the root CLI command suite.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(onboard.NewOnboardCmd())
	rootCmd.AddCommand(onboard.NewVerifyCmd())
}
