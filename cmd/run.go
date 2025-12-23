package cmd

import (
	"github.com/samama/firaaq/internal"
	"github.com/spf13/cobra"
)

// NewRunCommand implements and returns the run command.
func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short:                 "Run a command inside a new Container.",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(1),
		RunE:                  internal.Run,
	}

	attachRunFlags(cmd)
	return cmd
}

// attachRunFlags declares the run command's CLI options.
func attachRunFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringP("host", "", "", "Container Hostname")
	flags.IntP("memory", "m", 100, "Limit memory access in MB")
	flags.IntP("swap", "s", 20, "Limit swap access in MB")
	flags.Float64P("cpus", "c", 2, "Limit CPUs")
	flags.IntP("pids", "p", 128, "Limit number of processes")
	flags.BoolP("detach", "d", false, "run command in the background")
}