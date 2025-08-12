package cmd

import (
	"fmt"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var cleanAll bool

var CleanCmd = &cobra.Command{
	Use:   "clean [list1 names...]",
	Short: "Remove all completed tasks from the specified lists or just the current list if none are specified.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if cleanAll {
			cleanAllLists()
		} else if len(args) > 0 {
			cleanSpecifiedLists(args)
		} else {
			cleanCurrentList()
		}
	},
}

func setUpClean() {
	RootCmd.AddCommand(CleanCmd)
	CleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Remove all completed tasks from all lists.")
}

func cleanAllLists() {
	core.WithDefaultDB(func(db *core.DB) {
		err := db.CleanAllLists()
		if err != nil {
			core.Abort(fmt.Sprintf("Error cleaning all todo-lists: %v", err))
		}
		core.Success("Cleaned completed tasks from all todo-lists.")
	})
}

func cleanCurrentList() {
	core.WithDefaultDB(func(db *core.DB) {
		current, err := db.GetCurrentListName()
		if err != nil {
			core.Abort(fmt.Sprintf("Error accessing current todo-list: %v", err))
		}
		if current == "" {
			core.Abort("No current todo-list is set.")
		}

		err = db.CleanCurrentList()
		if err != nil {
			core.Abort(fmt.Sprintf("Error cleaning current todo-list: %v", err))
		}

		core.Success(fmt.Sprintf("Cleaned the following:\n%s", core.ListLists([]string{current}, "  ")))
	})
}

func cleanSpecifiedLists(names []string) {
	core.WithDefaultDB(func(db *core.DB) {
		err := db.CleanLists(names)
		if err != nil {
			core.Abort(fmt.Sprintf("Error cleaning specified todo-lists: %v", err))
		}
		core.Success(fmt.Sprintf("Cleaned the following:\n%s", core.ListLists(names, "  ")))
	})
}
