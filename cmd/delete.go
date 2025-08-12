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
	RunE: func(cmd *cobra.Command, args []string)  error{
		if deleteAll {
			return deleteAllLists()
		} else if len(args) == 0 {
			return deleteCurrentList()
		} else {
			return deleteSpecifiedLists(args)
		}
	},
}

func setUpDelete() {
	RootCmd.AddCommand(DeleteCmd)
	DeleteCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "Delete all lists")
}

func deleteAllLists() error {
	return core.WithDefaultDB(func(db *core.DB) error {
		err := db.DeleteAllLists()
		if err != nil {
			return fmt.Errorf("could not delete all todo-lists due to the following error\n\t %v", err)
		}
		core.Success("Deleted all todo-lists.")
		return nil
	})
}

func deleteCurrentList() error {
	return core.WithDefaultDB(func(db *core.DB) error {
		current, err := db.GetCurrentListName()
		if err != nil {
			return fmt.Errorf("could not access current todo-list due to the following error\n\t %v", err)
		}
		if current == "" {
			return fmt.Errorf("no current todo-list is set")
		}

		err = db.DeleteCurrentList()
		if err != nil {
			return fmt.Errorf("could not delete current todo-list due to the following error\n\t %v", err)
		}

		core.Success(fmt.Sprintf("Deleted the following:\n%s", core.ListLists([]string{current}, "  ")))
		return nil
	})
}

func deleteSpecifiedLists(names []string) error{
	return core.WithDefaultDB(func(db *core.DB) error {
		err := db.DeleteLists(names)
		if err != nil {
			return fmt.Errorf("could not delete specified todo-lists due to the following error\n\t %v", err)
		}
		core.Success(fmt.Sprintf("Deleted the following:\n%s", core.ListLists(names, "  ")))
		return nil
	})
}
