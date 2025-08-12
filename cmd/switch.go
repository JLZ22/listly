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
	Run: func(cmd *cobra.Command, args []string) {
		listName := args[0]
		core.WithDefaultDB(func(db *core.DB) {
			exists, err := db.ListExists(listName)
			if err != nil {
				core.Abort(fmt.Sprintf("Error checking if list %s exists: %v", listName, err))
			}
			if !exists {
				core.Abort(fmt.Sprintf("List %s does not exist. Cannot switch to it.", listName))
			}

			if err := db.SetCurrentListName(listName); err != nil {
				core.Abort(fmt.Sprintf("Error switching to list %s: %v", listName, err))
			}
		})
		core.Success(fmt.Sprintf("Switched to todo-list '%s'", listName))
	},
}

func setUpSwitch() {
	RootCmd.AddCommand(SwitchCmd)
}
