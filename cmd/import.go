package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type taskDTO struct {
	Description string `json:"description" yaml:"description"`
	Done        bool   `json:"done" yaml:"done"`
}

type listDTO struct {
	Title string `json:"title" yaml:"title"`
	Tasks []taskDTO
}

var ImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import tasks from a file. Supported formats: JSON, YAML",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := args[0]
		content, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		lists, err := fileToData(content, filepath.Ext(fileName))
		if err != nil {
			return err
		}
		err = core.WithDefaultDB(func(db *core.DB) error {
			for _, list := range lists {
				exists, err := db.ListExists(list.Info.Name)
				if err != nil {
					return err
				}
				if exists {
					return fmt.Errorf("failed to import because list %q already exists", list.Info.Name)
				}

				err = db.SaveList(list)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Printf("Imported the following lists:\n")
		for _, list := range lists {
			fmt.Println("  - ", list.Info.Name)
		}
		return nil
	},
}

func setUpImport() {
	RootCmd.AddCommand(ImportCmd)
}

func fileToData(content []byte, ext string) ([]core.List, error) {
	var dtos []listDTO
	var lists []core.List
	var err error

	// unmarshal based on file extension
	switch ext {
	case ".json":
		dec := json.NewDecoder(bytes.NewReader(content))
		dec.DisallowUnknownFields() // ensure no unknown fields
		err = dec.Decode(&dtos)
	case ".yaml":
		dec := yaml.NewDecoder(bytes.NewReader(content))
		dec.KnownFields(true)
		err = dec.Decode(&dtos)
	default:
		return lists, fmt.Errorf("unsupported file format: \"%s\". Supported formats are JSON and YAML", ext)
	}
	if err != nil {
		return lists, err
	}

	// convert dtos into lists 
	lists = make([]core.List, len(dtos))
	for i, dto := range dtos {
		list := core.NewList(dto.Title)
		for _, task := range dto.Tasks {
			list.AddNewTask(task.Description, task.Done)
		}
		lists[i] = list
	}
	return lists, nil
}
