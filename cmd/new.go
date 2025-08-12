package cmd

import (
	"fmt"
	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var NewCmd = &cobra.Command{
	Use:   "new [list name]",
	Short: "Create a new todo list.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		core.WithDefaultDB(func(db *core.DB) {
			for _, listName := range args {
				createNewList(db, listName)
			}
			core.Success(fmt.Sprintf("Created new todo-lists:\n%s", core.ListLists(args, "  ")))
		})
	},
}

func setUpNew() {
	RootCmd.AddCommand(NewCmd)
}

func createNewList(db *core.DB, listName string) {
	exists, err := db.ListExists(listName)
	if err != nil {
		core.Abort(fmt.Sprintf("Error accessing the database: %v", err))
	}
	if exists {
		core.Abort(fmt.Sprintf("List \"%s\" already exists.", listName))
	}

	newList := core.NewList(listName)
	err = db.SaveList(newList)
	if err != nil {
		core.Abort(fmt.Sprintf("Error saving new list: %v", err))
	}
	err = db.SetCurrentListName(listName)
	if err != nil {
		core.Abort(fmt.Sprintf("Error setting current list name: %v", err))
	}
}
