package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlz22/listly/core"
	"github.com/jlz22/listly/tui"
	"github.com/spf13/cobra"
)

var OpenCmd = &cobra.Command{
	Use:   "open [list name]",
	Short: "Open the TUI for the specified list or the current list if no list is specified.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.WithDefaultDB(
			func(db *core.DB) error {
				var listName string
				if len(args) > 0 {
					listName = args[0]
				} else {
					var err error
					listName, err = db.GetCurrentListName()
					if err != nil {
						return err
					}
					if listName == "" {
						return fmt.Errorf("no list chosen. specify a list name or use\n\n\tlistly switch <list name>\n ")
					}
				}

				// switch to the specified list
				err := db.SetCurrentListName(listName)
				if err != nil {
					return err
				}

				exists, err := db.ListExists(listName)
				if err != nil {
					return err
				}
				if !exists {
					return fmt.Errorf("list %q does not exist", listName)
				}

				m, err := tui.NewModel(db, listName)
				if err != nil {
					return err
				}

				tuiProgram := tea.NewProgram(m)
				_, err = tuiProgram.Run()
				if err != nil {
					return err
				}
				fmt.Print("\033[H\033[2J") // clear the screen
				return nil
			},
		)
	},
}

func setUpOpen() {
	RootCmd.AddCommand(OpenCmd)
}
