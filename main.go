package main

import (
	"github.com/samama/firaaq/cmd"
)

func main() {
	// Build the root CLI and attach the supported container commands before executing.
	rootCmd := cmd.NewFiraaqCommand()
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewForkCommand())
	rootCmd.AddCommand(cmd.NewExecCommand())
	rootCmd.AddCommand(cmd.NewPsCommand())
	rootCmd.AddCommand(cmd.NewImagesCommand())
	rootCmd.Execute()
}
