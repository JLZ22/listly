package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [list name]",
	Short: "Print all tasks in the current or specified list.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var listName string
		if len(args) > 0 {
			listName = args[0]
		} else {
			listName = "current list" // replace with curr list retrieval
		}
		fmt.Printf("Showing tasks for list: %s\n", listName)
	},
}

func setUpShow() {
	RootCmd.AddCommand(showCmd)
}
