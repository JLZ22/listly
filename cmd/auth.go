package cmd

import (
	"github.com/spf13/cobra"
	"github.com/jlz22/listly/core"
)

var AuthCmd = &cobra.Command{
	Use:   "auth <api key>",
	Short: "Add Google Gemini API key.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := args[0]
		core.WithDefaultDB(func(db *core.DB) error {
			return db.SetAPIKey(apiKey)
		})
		return nil
	},
}

func setUpAuth() {
	RootCmd.AddCommand(AuthCmd)
}