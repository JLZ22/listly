package cmd

import (
	"fmt"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var deleteAll bool

var DeleteCmd = &cobra.Command{
	Use:   "delete <list name> [more list names...]",
	Short: "Delete the specified list(s).",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteAll {
			return deleteAllLists()
		} else if len(args) == 0 {
			return cmd.Help()
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

func deleteSpecifiedLists(names []string) error {
	return core.WithDefaultDB(func(db *core.DB) error {
		// separate the lists that exist vs the ones that don't
		found := []string{}
		notFound := []string{}
		seen := map[string]struct{}{}
		for _, name := range names {
			_, ok := seen[name]
			if ok {
				continue
			}
			seen[name] = struct{}{}

			exists, err := db.ListExists(name)
			if err != nil {
				return fmt.Errorf("could not check if list %s exists due to the following error\n\t %v", name, err)
			}

			if exists {
				found = append(found, name)
			} else {
				notFound = append(notFound, name)
			}
		}

		// Delete the found lists
		err := db.DeleteLists(found)
		if err != nil {
			return fmt.Errorf("could not delete specified todo-lists due to the following error\n\t %v", err)
		}

		core.Success(fmt.Sprintf("Deleted the following:\n%s", core.ListLists(found, "  ")))
		if len(notFound) > 0 {
			fmt.Printf("Could not find the following:\n%s", core.ListLists(notFound, "  "))
		}
		return nil
	})
}
