package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/Siddhant-K-code/sentinel-ai/internal/engine"
	"github.com/Siddhant-K-code/sentinel-ai/internal/policy"
)

func prCmd() *cobra.Command {
	var (
		title     string
		body      string
		draft     bool
		planPath  string
		policyPath string
	)

	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Create a pull request with proposed changes",
		Long: `Create a pull request containing the changes from a plan file.
Requires GitHub CLI (gh) to be installed and authenticated.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Load policy
			pol, err := policy.Load(policyPath)
			if err != nil {
				return err
			}

			// Create engine
			e, err := engine.New(ctx, engine.Options{
				Policy: pol,
			})
			if err != nil {
				return err
			}

			// Create PR
			result, err := e.CreatePR(ctx, engine.PROptions{
				Title:    title,
				Body:     body,
				Draft:    draft,
				PlanPath: planPath,
			})
			if err != nil {
				return err
			}

			cmd.Printf("Pull request created: %s\n", result.URL)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "PR title (required)")
	cmd.Flags().StringVar(&body, "body", "", "PR body (file path or text)")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft PR")
	cmd.Flags().StringVar(&planPath, "plan", "", "Path to plan file (required)")
	cmd.Flags().StringVar(&policyPath, "policy", "./.sentinel/policy.yaml", "Path to policy file")

	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("plan")

	return cmd
}
