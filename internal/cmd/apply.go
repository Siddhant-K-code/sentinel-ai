package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/Siddhant-K-code/sentinel-ai/internal/engine"
	"github.com/Siddhant-K-code/sentinel-ai/internal/policy"
)

func applyCmd() *cobra.Command {
	var (
		planPath     string
		approveLevel string
		policyPath   string
	)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply patches from a plan file",
		Long: `Apply patches generated from a previous scan.
Requires explicit approval level and reads from a plan file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Load policy
			pol, err := policy.Load(policyPath)
			if err != nil {
				return err
			}

			// Read plan file
			planData, err := os.ReadFile(planPath)
			if err != nil {
				return err
			}

			var plan engine.Plan
			if err := json.Unmarshal(planData, &plan); err != nil {
				return err
			}

			// Create engine
			e, err := engine.New(ctx, engine.Options{
				Policy: pol,
			})
			if err != nil {
				return err
			}

			// Apply patches
			result, err := e.Apply(ctx, plan, approveLevel)
			if err != nil {
				return err
			}

			// Output results
			if result.Success {
				cmd.Println("Patches applied successfully")
			} else {
				cmd.Printf("Patch application failed: %s\n", result.Error)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&planPath, "plan", "", "Path to plan file (required)")
	cmd.Flags().StringVar(&approveLevel, "approve-level", "low", "Approval level (low, medium, high)")
	cmd.Flags().StringVar(&policyPath, "policy", "./.sentinel/policy.yaml", "Path to policy file")

	_ = cmd.MarkFlagRequired("plan")

	return cmd
}
