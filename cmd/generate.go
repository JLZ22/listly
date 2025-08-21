package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jlz22/listly/core"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

var instructions = "You are an experienced developer who needs to make a " +
	"set of todo lists for your junior developer so that they can complete whatever is " +
	"described below. Each list should break the tasks into clear steps as a developer " +
	"roadmap. Use UpperCamelCase for list titles."

var timeoutFlagValue int

var GenerateCmd = &cobra.Command{
	Use:   "generate <file>",
	Short: "Generate lists based on descriptions in a txt file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// get text content
		fileName := args[0]
		content, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}

		// get content using Gemini
		return core.WithDefaultDB(func(db *core.DB) error {
			// get API key
			apiKey, err := db.GetAPIKey()
			if err != nil {
				return err
			}

			// create client
			ctx := context.Background()
			client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
			if err != nil {
				return err
			}

			// Update the prompt to exclude names already in the database
			allInfo, err := db.GetInfo()
			if err != nil {
				return err
			}
			existingLists := ""
			for _, info := range allInfo {
				existingLists += "\"" + info.Name + "\", "
			}
			instructions += " Exclude the following names from the lists: " + strings.TrimSuffix(existingLists, ", ") + "."

			// generate lists
			result, err := generateLists(ctx, client, content)
			if err != nil {
				return err
			}

			// convert Gemini output to List type
			lists, err := fileToData([]byte(result.Text()), ".json")
			if err != nil {
				return err
			}
			for _, list := range lists {
				fmt.Println()
				fmt.Printf("%v", &list)
			}

			// confirm changes
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nAdd these to your lists? (y/n) ")
			for {
				char, _, err := reader.ReadRune()
				if err != nil {
					return err
				}
				if char == 'y' || char == 'Y' {
					// save lists to the DB
					for _, list := range lists {
						db.SaveList(list)
					}
					break
				} else if char == 'n' || char == 'N' {
					// discard changes
					fmt.Println("Lists discarded.")
					break
				}
			}
			fmt.Println("Success! All lists added.")

			return nil
		})
	},
}

func setUpGenerate() {
	RootCmd.AddCommand(GenerateCmd)
	GenerateCmd.Flags().IntVarP(&timeoutFlagValue, "timeout", "t", 120, "Number of seconds to wait before quitting")
}

// Generate a set of lists using Gemini with the given content. Print a spinner and handle timeout while working.
func generateLists(ctx context.Context, client *genai.Client, content []byte) (*genai.GenerateContentResponse, error) {
	functionCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	timeoutTimer := time.NewTimer(time.Duration(timeoutFlagValue)*time.Second + time.Millisecond*300)

	type resultStruct struct {
		response *genai.GenerateContentResponse
		err      error
	}
	resultCh := make(chan resultStruct, 1)

	// print generating status with updating timer
	go func() {
		secondTicker := time.NewTicker(time.Duration(1) * time.Second)
		defer secondTicker.Stop()
		start := time.Now()

		for {
			select {
			case <-timeoutTimer.C:
				return
			case <-functionCtx.Done():
				return
			case <-secondTicker.C:
				elapsed := int(time.Since(start).Seconds())
				fmt.Printf("\rGenerating todo lists... %2d / %ds elapsed", elapsed, timeoutFlagValue)
			}
		}
	}()

	// goroutine to generate lists with Gemini
	go func() {
		result, err := client.Models.GenerateContent(
			functionCtx,
			"gemini-2.5-flash",
			genai.Text(instructions+"\n\n"+string(content)),
			core.GeminiConfig,
		)
		resultCh <- resultStruct{result, err}
	}()

	// wait for result or timeout
	select {
	case <-timeoutTimer.C:
		fmt.Println()
		return nil, fmt.Errorf("timed out after %d seconds", timeoutFlagValue)
	case result := <-resultCh:
		fmt.Println()
		if result.err != nil {
			return nil, fmt.Errorf("gemini error - %v", result.err)
		}
		fmt.Println("Success!")
		return result.response, nil
	}
}
