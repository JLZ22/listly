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
	Use:   "export <file> [list names...] ",
	Short: "Export list to a file. Uses current list if no list name(s) provided. Supported formats: JSON, YAML",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lists := make([]core.List, max(1, len(args) - 1))
		err := core.WithDefaultDB(func(db *core.DB) error {
			var fileName string

			if len(args) == 1 { // no list name specified so use the current list
				listName, err := db.GetCurrentListName()
				if err != nil {
					return err
				}
				list, err := db.GetList(listName)
				if err != nil {
					return err
				}
				fileName = args[0]
				lists[0] = list
			} else {
				fileName = args[0]
				for i := 1 ; i < len(args) ; i++ {
					list, err := db.GetList(args[i])
					if err != nil {
						return err
					}
					lists[i-1] = list
				}
			}

			content, err := dataToFile(lists, filepath.Ext(fileName))
			if err != nil {
				return err
			}
			return os.WriteFile(fileName, content, 0644)
		})
		if err != nil {
			return err
		}
		fmt.Printf("Exported the following lists to \"%s\":\n", args[0])
		for _, list := range lists {
			fmt.Println("  - ", list.Info.Name)
		}
		return nil
	},
}

func setUpExport() {
	RootCmd.AddCommand(ExportCmd)
}

func dataToFile(lists []core.List, ext string) ([]byte, error) {
	var content []byte
	var err error
	var dtos []listDTO = make([]listDTO, len(lists))

	// convert lists of lists to list of data transfer objects (listDTO)
	for i, list := range lists {
		var dto listDTO
		dto.Title = list.Info.Name
		for _, task := range list.Tasks {
			dto.Tasks = append(dto.Tasks, taskDTO{
				Description: task.Description,
				Done:        task.Done,
			})
		}
		dtos[i] = dto
	}

	// marshal every DTO based on file extension
	switch ext {
	case ".json":
		content, err = json.MarshalIndent(dtos, "", "  ")
	case ".yaml":
		content, err = yaml.Marshal(dtos)
	default:
		return content, fmt.Errorf("unsupported file format: \"%s\". Supported formats are JSON and YAML", ext)
	}
	if err != nil {
		return content, err
	}
	return content, nil
}
