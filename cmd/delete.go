package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [list name]",
	Short: "Delete the specified list.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listName := args[0]
		// call db to delete the list
		fmt.Printf("Deleting list: %s", listName)
	},
}

func setUpDelete() {
	RootCmd.AddCommand(deleteCmd)
}
