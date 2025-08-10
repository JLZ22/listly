package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [list name]",
	Short: "Switch to the specified todo list.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listName := args[0]
		// call db to switch the current list
		fmt.Printf("Switching to list: %s\n", listName)
	},
}

func setUpSwitch() {
	RootCmd.AddCommand(switchCmd)
}
