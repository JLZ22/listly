package cmd

import (
	"fmt"
	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show [list name]",
	Short: "Print all tasks in the current or specified list.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.WithDefaultDB(func(db *core.DB) error {
			// Get the list name (current or specified)
			var listName string
			if len(args) > 0 {
				listName = args[0]
			} else {
				var err error
				listName, err = db.GetCurrentListName()
				if err != nil {
					return fmt.Errorf("could not retrieve current list name due to the following error\n\t %v", err)
				}
				if listName == "" {
					return fmt.Errorf("no list selected")
				}
			}

			// Get the list struct
			list, err := db.GetList(listName)
			if err != nil {
				return fmt.Errorf("could not retrieve list %s due to the following error\n\t %v", listName, err)
			}

			// Print the tasks in the list
			if len(list.Tasks) == 0 {
				fmt.Printf("No tasks found in list '%s'\n", listName)
				return nil
			}
			completed, pending := core.SplitByCompletion(list)
			fmt.Printf("%s\n\n", listName)
			for _, task := range pending {
				fmt.Printf("- [ ] %s\n", task.Description)
			}
			for _, task := range completed {
				fmt.Printf("- [x] %s\n", task.Description)
			}
			return nil
		})
	},
}

func setUpShow() {
	RootCmd.AddCommand(ShowCmd)
}
