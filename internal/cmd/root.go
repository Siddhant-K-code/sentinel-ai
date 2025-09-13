package cmd

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "sentinel-ai",
		Short: "Security scanning and dead-code detection CLI (POC)",
		Long: `Sentinel-AI is a hermetic CLI for security scanning and dead-code detection.
It performs SAST analysis, identifies dead code, and can propose or apply patches.

⚠️  WARNING: This is a proof of concept (POC) for personal use only.
Not intended for production environments.`,
	}

	root.AddCommand(scanCmd())
	root.AddCommand(applyCmd())
	root.AddCommand(prCmd())

	return root
}
