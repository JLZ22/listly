package main

import (
	cmd "github.com/jlz22/listly/cmd"
)

func main() {
	cmd.SetUp()
	cmd.RootCmd.Execute()
}
