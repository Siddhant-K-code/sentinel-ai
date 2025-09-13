package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/Siddhant-K-code/sentinel-ai/internal/engine"
	"github.com/Siddhant-K-code/sentinel-ai/internal/policy"
)

func scanCmd() *cobra.Command {
	var (
		repo       string
		agentPath  string
		policyPath string
		sarifOut   string
		planOut    string
		logOut     string
		doSec      bool
		doDead     bool
	)

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan repository for security issues and dead code",
		Long: `Perform security scanning and dead-code detection on the specified repository.
Produces SARIF output and patch plans without modifying the codebase.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Load policy
			pol, err := policy.Load(policyPath)
			if err != nil {
				return err
			}

			// Create engine
			e, err := engine.New(ctx, engine.Options{
				Repo:      repo,
				AgentPath: agentPath,
				Policy:    pol,
				LogPath:   logOut,
			})
			if err != nil {
				return err
			}

			// Run scan
			res, err := e.Scan(ctx, engine.ScanOpts{
				Security: doSec,
				DeadCode: doDead,
			})
			if err != nil {
				return err
			}

			// Write outputs
			if sarifOut != "" {
				if err := os.WriteFile(sarifOut, res.SARIF, 0644); err != nil {
					return err
				}
			}

			if planOut != "" {
				b, err := json.MarshalIndent(res.Plan, "", "  ")
				if err != nil {
					return err
				}
				if err := os.WriteFile(planOut, b, 0644); err != nil {
					return err
				}
			}

			os.Exit(res.ExitCode)
			return nil
		},
	}

	cmd.Flags().StringVar(&repo, "repo", ".", "Repository path to scan")
	cmd.Flags().StringVar(&agentPath, "agent", "./AGENT.md", "Path to AGENT.md file")
	cmd.Flags().StringVar(&policyPath, "policy", "./.sentinel/policy.yaml", "Path to policy file")
	cmd.Flags().StringVar(&sarifOut, "sarif", "", "SARIF output file path")
	cmd.Flags().StringVar(&planOut, "plan", "", "Plan output file path")
	cmd.Flags().StringVar(&logOut, "log", "", "Log output file path")
	cmd.Flags().BoolVar(&doSec, "security", false, "Enable security scanning")
	cmd.Flags().BoolVar(&doDead, "dead-code", false, "Enable dead-code detection")

	return cmd
}
