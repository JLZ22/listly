package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var KmapCmd = &cobra.Command{
	Use:   "kmap [command]",
	Short: "Manage key-mappings for the TUI",
}

var KmapShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show which file is being used for the current keymap.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// get the file path
		var pth string
		err := core.WithDefaultDB(func(db *core.DB) error {
			var err error
			pth, err = db.GetKmapPath()
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		if pth == "" { // No file set.
			fmt.Println("No file set. Using defaults.")
		} else { // File is set.
			fmt.Printf("Using %s for key-mapping.\n", pth)
			// File doesn't exist.
			_, err = os.Stat(pth)
			if err != nil {
				fmt.Println("WARNING: File does not exist. Using defaults.")
			}
		}

		return nil
	},
}

var KmapSetCmd = &cobra.Command{
	Use:   "set <file>",
	Short: "Sets the file to refer to for key-mappings. Refer to the README for file formatting.",
	Long:  "Sets the file to refer to for key-mappings. Fails if file does not exist. Refer to the README for file formatting.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pth := args[0]
		absPath, err := filepath.Abs(pth)
		if err != nil {
			return err
		}

		err = core.WithDefaultDB(func(db *core.DB) error {
			return db.SetKmapPath(absPath)
		})
		if err != nil {
			return err
		}

		fmt.Printf("Success! Refering to %s for key-mappings.\n", pth)
		return nil
	},
}

var KmapClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Revert to default bindings. Does not edit files.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.WithDefaultDB(func(db *core.DB) error {
			return db.SetKmapPath("")
		})
		if err != nil {
			return err
		}

		fmt.Println("Success! Key-bindings reset to default.")
		return nil
	},
}

func setUpKmap() {
	RootCmd.AddCommand(KmapCmd)
	KmapCmd.AddCommand(KmapShowCmd)
	KmapCmd.AddCommand(KmapSetCmd)
	KmapCmd.AddCommand(KmapClearCmd)
}
