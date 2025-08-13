package cmd

import (
	"fmt"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var SwitchCmd = &cobra.Command{
	Use:   "switch [list name]",
	Short: "Switch to the specified todo list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listName := args[0]
		return core.WithDefaultDB(func(db *core.DB) error {
			exists, err := db.ListExists(listName)
			if err != nil {
				return fmt.Errorf("could not check if list %s exists due to the following error\n\t %v", listName, err)
			}

			if !exists {
				return fmt.Errorf("list %s does not exist - cannot switch to it", listName)
			}

			if err := db.SetCurrentListName(listName); err != nil {
				return fmt.Errorf("could not switch to list %s due to the following error\n\t %v", listName, err)
			}
			core.Success(fmt.Sprintf("Switched to todo-list '%s'", listName))
			return nil
		})
	},
}

func setUpSwitch() {
	RootCmd.AddCommand(SwitchCmd)
}
