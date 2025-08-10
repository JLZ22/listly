package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var cleanAll bool

var cleanCmd = &cobra.Command{
	Use:   "clean [list1 names...]",
	Short: "Remove all completed tasks from the specified lists or just the current list if none are specified.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if cleanAll {
			fmt.Println("Cleaning all lists...")
		} else if len(args) > 0 {
			fmt.Printf("Cleaning specified lists: %v\n", args)
		} else {
			fmt.Println("Cleaning current list...")
		}
	},
}

func setUpClean() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Remove all completed tasks from all lists.")
}
