package main

import (
	"log"

	cmd "github.com/jlz22/listly/cmd"
)

func main() {
	cmd.SetUp()
	err := cmd.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
