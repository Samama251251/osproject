package main

import (
	"github.com/samama/firaaq/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Build the root CLI and attach the supported container commands before executing.
	rootCmd := cmd.NewFiraaqCommand()
	commandBuilders := []func() *cobra.Command{
		cmd.NewRunCommand,
		cmd.NewForkCommand,
		cmd.NewExecCommand,
		cmd.NewPsCommand,
		cmd.NewImagesCommand,
	}
	for _, build := range commandBuilders {
		rootCmd.AddCommand(build())
	}
	rootCmd.Execute()
}