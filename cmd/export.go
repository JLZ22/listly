package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ExportCmd = &cobra.Command{
	Use:   "export [list name] <file>",
	Short: "Export list to a file. Uses current list if no list name provided. Supported formats: JSON, YAML",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.WithDefaultDB(func(db *core.DB) error {
			var listName string
			var fileName string

			if len(args) == 0 {
				return cmd.Help()
			} else if len(args) == 1 { // no list name specified so use the current list
				var err error
				listName, err = db.GetCurrentListName()
				if err != nil {
					return err
				}
				fileName = args[0]
			} else {
				listName = args[0]
				fileName = args[1]
			}

			list, err := db.GetList(listName)
			if err != nil {
				return err
			}
			return dataToFile(list, fileName)
		})
	},
}

func setUpExport() {
	RootCmd.AddCommand(ExportCmd)
}

func dataToFile(list core.List, fileName string) error {
	var content []byte
	var err error
	var dto listDTO

	// convert list to data transfer object (listDTO)
	dto.Title = list.Info.Name
	for _, task := range list.Tasks {
		dto.Tasks = append(dto.Tasks, taskDTO{
			Description: task.Description,
			Done:        task.Done,
		})
	}

	// marshal based on file extension
	ext := filepath.Ext(fileName)

	switch ext {
	case ".json":
		content, err = json.MarshalIndent(dto, "", "  ")
	case ".yaml":
		content, err = yaml.Marshal(dto)
	default:
		return fmt.Errorf("unsupported file format: \"%s\". Supported formats are JSON and YAML", ext)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, content, 0644)
}