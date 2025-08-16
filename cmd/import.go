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
		list, err := fileToData(content, filepath.Ext(fileName))
		if err != nil {
			return err
		}
		return core.WithDefaultDB(func(db *core.DB) error {
			exists, err := db.ListExists(list.Info.Name)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("failed to import because list %q already exists", list.Info.Name)
			}
			return db.SaveList(list)
		})
	},
}

func setUpImport() {
	RootCmd.AddCommand(ImportCmd)
}

func fileToData(content []byte, ext string) (core.List, error) {
	var dto listDTO
	var data core.List
	var err error

	// unmarshal based on file extension
	switch ext {
	case ".json":
		dec := json.NewDecoder(bytes.NewReader(content))
		dec.DisallowUnknownFields() // ensure no unknown fields
		err = dec.Decode(&dto)
	case ".yaml":
		dec := yaml.NewDecoder(bytes.NewReader(content))
		dec.KnownFields(true)
		err = dec.Decode(&dto)
	default:
		return data, fmt.Errorf("unsupported file format: \"%s\". Supported formats are JSON and YAML", ext)
	}
	if err != nil {
		return data, err
	}

	data = core.NewList(dto.Title)
	for _, task := range dto.Tasks {
		data.AddNewTask(task.Description, task.Done)
	}
	return data, nil
}
