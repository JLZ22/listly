package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var OpenCmd = &cobra.Command{
	Use:   "open [list name]",
	Short: "Open the TUI for the specified list or the current list if no list is specified.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// DONT FORGET THE LIST NOT FOUND CASE
		var listName string
		if len(args) > 0 {
			listName = args[0]
		} else {
			listName = "current list" // replace with curr list retrieval
		}
		fmt.Printf("Opening TUI for list: %s\n", listName) // delete later
		fmt.Printf("\n\tTODO NOT IMPLEMENTED\n\n")
	},
}

func setUpOpen() {
	RootCmd.AddCommand(OpenCmd)
}
