package cmd

import (
	"fmt"
	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var RenameCmd = &cobra.Command{
	Use:   "rename [old list name] [new list name]",
	Short: "Rename an existing todo list.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]

		core.WithDefaultDB(func(db *core.DB) {
			err := db.RenameList(oldName, newName)
			if err != nil {
				core.Abort(fmt.Sprintf("Error renaming todo-list: %v", err))
			}
			core.Success(fmt.Sprintf("Renamed todo-list '%s' to '%s'", oldName, newName))
		})
	},
}

func setUpRename() {
	RootCmd.AddCommand(RenameCmd)
}
