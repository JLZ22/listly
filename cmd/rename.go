package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [old list name] [new list name]",
	Short: "Rename an existing todo list.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]
		// call db to rename the list
		fmt.Printf("Renaming list '%s' to '%s'\n", oldName, newName)
	},
}

func setUpRename() {
	RootCmd.AddCommand(renameCmd)
}
