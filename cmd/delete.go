package cmd

import (
	"fmt"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var deleteAll bool

var DeleteCmd = &cobra.Command{
	Use:   "delete [list name]",
	Short: "Delete the specified list.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if deleteAll {
			deleteAllLists()
		} else if len(args) == 0 {
			deleteCurrentList()
		} else {
			deleteSpecifiedLists(args)
		}
	},
}

func setUpDelete() {
	RootCmd.AddCommand(DeleteCmd)
	DeleteCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "Delete all lists")
}

func deleteAllLists() {
	core.WithDefaultDB(func(db *core.DB) {
		err := db.DeleteAllLists()
		if err != nil {
			core.Abort(fmt.Sprintf("Error deleting all todo-lists: %v", err))
		}
		core.Success("Deleted all todo-lists.")
	})
}

func deleteCurrentList() {
	core.WithDefaultDB(func(db *core.DB) {
		current, err := db.GetCurrentListName()
		if err != nil {
			core.Abort(fmt.Sprintf("Error accessing current todo-list: %v", err))
		}
		if current == "" {
			core.Abort("No current todo-list is set.")
		}

		err = db.DeleteCurrentList()
		if err != nil {
			core.Abort(fmt.Sprintf("Error deleting current todo-list: %v", err))
		}

		core.Success(fmt.Sprintf("Deleted the following:\n%s", core.ListLists([]string{current}, "  ")))
	})
}

func deleteSpecifiedLists(names []string) {
	core.WithDefaultDB(func(db *core.DB) {
		err := db.DeleteLists(names)
		if err != nil {
			core.Abort(fmt.Sprintf("Error deleting specified todo-lists: %v", err))
		}
		core.Success(fmt.Sprintf("Deleted the following:\n%s", core.ListLists(names, "  ")))
	})
}
