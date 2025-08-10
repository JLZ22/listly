package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [list name]",
	Short: "Create a new todo list.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listName := args[0]
		// call db to create the list
		fmt.Printf("Creating new list: %s\n", listName)
	},
}

func setUpNew() {
	RootCmd.AddCommand(newCmd)
}
