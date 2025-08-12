package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
)

var useQuotes bool

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display the names of all todo lists along with their task counts.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		core.WithDefaultDB(func(db *core.DB) {
			allInfo, err := db.GetInfo()
			if err != nil {
				core.Abort(fmt.Sprintf("Failed to retrieve lists: %v", err))
			}

			// end early if no lists found
			if len(allInfo) == 0 {
				core.Abort("No lists found. Create a new one with\n\n\tlistly new <list name> ...")
			}

			// Find the maximum length of list names for formatting
			currentListName, err := db.GetCurrentListName()
			if err != nil {
				core.Abort(fmt.Sprintf("Failed to retrieve current list name: %v", err))
			}
			maxLen := 12
			for _, info := range allInfo {
				name := info.Name
				nameLen := len(name)
				if name == currentListName {
					nameLen += len(" (current)") + 2
				}
				if nameLen > maxLen {
					maxLen = nameLen
				}
			}

			// Convert map to slice for sorting
			var allInfoSlice []core.ListInfo
			for _, info := range allInfo {
				allInfoSlice = append(allInfoSlice, info)
			}

			// sort before displaying
			sort.Slice(allInfoSlice, func(i, j int) bool {
				return allInfoSlice[i].Name < allInfoSlice[j].Name
			})

			// display in a table format
			printRow("List Name", "Pending", "Done", "Total", maxLen)
			fmt.Println(strings.Repeat("=", maxLen+21)) // Print a separator line
			for _, info := range allInfoSlice {
				var name string
				if useQuotes {
					name = fmt.Sprintf("\"%s\"", info.Name)
				} else {
					name = info.Name
				}
				if info.Name == currentListName {
					printRow(
						fmt.Sprintf("%s (current)", name),
						strconv.Itoa(info.NumPending),
						strconv.Itoa(info.NumDone),
						strconv.Itoa(info.NumTasks),
						maxLen,
					)
				} else {
					printRow(
						name,
						strconv.Itoa(info.NumPending),
						strconv.Itoa(info.NumDone),
						strconv.Itoa(info.NumTasks),
						maxLen,
					)
				}
			}
		})
	},
}

func setUpList() {
	RootCmd.AddCommand(ListCmd)
	ListCmd.Flags().BoolVarP(&useQuotes, "quotes", "q", false, "Use quotes around list names")
}

func printRow(name, pending, done, total string, maxLen int) {
	fmt.Printf("%-*s %-*s %-*s %-*s\n", maxLen, name, 8, pending, 5, done, 0, total)
}
