package cmd

import (
	"fmt"
	"os"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Add Google Gemini API key.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Enter your Google Gemini API key: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		apiKey := string(password)
		err = core.WithDefaultDB(func(db *core.DB) error {
			return db.SetAPIKey(apiKey)
		})
		if err != nil {
			return err
		}
		fmt.Println("\nSuccessfully set API key.")
		return nil
	},
}

var AuthDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete your Google Gemini API key.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.WithDefaultDB(func(db *core.DB) error {
			return db.SetAPIKey("")
		})
		if err != nil {
			return err
		}
		fmt.Println("Successfully deleted API key.")
		return nil
	},
}

func setUpAuth() {
	RootCmd.AddCommand(AuthCmd)
	AuthCmd.AddCommand(AuthDeleteCmd)
}
