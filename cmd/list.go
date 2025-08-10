package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display the names of all todo lists along with their task counts.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// also include some indicator of what is the current list
		fmt.Println("Listing all todo lists...")
	},
}

func setUpList() {
	RootCmd.AddCommand(listCmd)
}
