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
	RunE: func(cmd *cobra.Command, args []string) error {
		oldName := args[0]
		newName := args[1]

		return core.WithDefaultDB(func(db *core.DB) error {
			err := db.RenameList(oldName, newName)
			if err != nil {
				return fmt.Errorf("could not rename todo-list due to the following error\n\t %v", err)
			}
			core.Success(fmt.Sprintf("Renamed todo-list '%s' to '%s'", oldName, newName))
			return nil
		})
	},
}

func setUpRename() {
	RootCmd.AddCommand(RenameCmd)
}
