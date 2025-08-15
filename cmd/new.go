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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Remove duplicates from args
		noDup := []string{}
		seen := map[string]struct{}{}
		for _, listName := range args {
			_, ok := seen[listName]
			if ok {
				continue
			}
			seen[listName] = struct{}{}
			noDup = append(noDup, listName)
		}

		return core.WithDefaultDB(func(db *core.DB) error {
			for _, listName := range noDup {
				exists, err := db.ListExists(listName)
				if err != nil {
					return err
				}
				if exists {
					return fmt.Errorf("list \"%s\" already exists. No new list created", listName)
				}
			}
			for _, listName := range noDup {
				err := createNewList(db, listName)
				if err != nil {
					return err
				}
			}
			core.Success(fmt.Sprintf("Created new todo-lists:\n%s", core.ListLists(noDup, "  ")))
			return nil
		})
	},
}

func setUpNew() {
	RootCmd.AddCommand(NewCmd)
}

func createNewList(db *core.DB, listName string) error {
	exists, err := db.ListExists(listName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("list \"%s\" already exists", listName)
	}

	newList := core.NewList(listName)
	err = db.SaveList(newList)
	if err != nil {
		return fmt.Errorf("could not save new list due to the following\n\t %v", err)
	}
	err = db.SetCurrentListName(listName)
	if err != nil {
		return fmt.Errorf("could not set current list name due to the following\n\t %v", err)
	}
	return nil
}
