package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "listly",
	Short: "Listly is a CLI task manager",
	Long: `
Listly is a task manager that lets you efficiently create and 
manage different todo lists with CLI commands. It also provides
a TUI that allows you to add/remove/edit tasks using Vim-style keybindings.`,
}

// can also be done with multiple init() functions, but it's easier to follow this way.
func SetUp() {
	setUpClean()
	setUpDelete()
	setUpList()
	setUpNew()
	setUpOpen()
	setUpRename()
	setUpShow()
	setUpSwitch()
	setUpImport()
	setUpExport()
}
