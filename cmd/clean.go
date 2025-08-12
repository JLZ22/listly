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
	RunE: func(cmd *cobra.Command, args []string) error {
		if cleanAll {
			return cleanAllLists()
		} else if len(args) > 0 {
			return cleanSpecifiedLists(args)
		} else {
			return cleanCurrentList()
		}
	},
}

func setUpClean() {
	RootCmd.AddCommand(CleanCmd)
	CleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Remove all completed tasks from all lists.")
}

func cleanAllLists() error {
	return core.WithDefaultDB(func(db *core.DB) error {
		err := db.CleanAllLists()
		if err != nil {
			return fmt.Errorf(" cleaning all todo-lists: %v", err)
		}
		core.Success("Cleaned completed tasks from all todo-lists.")
		return nil
	})
}

func cleanCurrentList() error {
	return core.WithDefaultDB(func(db *core.DB) error {
		current, err := db.GetCurrentListName()
		if err != nil {
			return fmt.Errorf("could not access current todo-list due to the following error\n\t %v", err)
		}
		if current == "" {
			return fmt.Errorf("no current todo-list is set")
		}

		err = db.CleanCurrentList()
		if err != nil {
			return fmt.Errorf("could not clean current todo-list due to the following error\n\t %v", err)
		}

		core.Success(fmt.Sprintf("Cleaned the following:\n%s", core.ListLists([]string{current}, "  ")))
		return nil
	})
}

func cleanSpecifiedLists(names []string) error {
	return core.WithDefaultDB(func(db *core.DB) error {
		err := db.CleanLists(names)
		if err != nil {
			return fmt.Errorf("could not clean specified todo-lists due to the following error\n\t %v", err)
		}
		core.Success(fmt.Sprintf("Cleaned the following:\n%s", core.ListLists(names, "  ")))
		return nil
	})
}
