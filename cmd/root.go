package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

const (
	layersPath     = "/var/lib/firaaq/images/layers"
	containersPath = "/var/run/firaaq/containers"
	netnsPath      = "/var/run/firaaq/netns"
)

var ErrNotPermitted = errors.New("operation not permitted")

// Make firaaq directories first.
func init() {
	os.MkdirAll(netnsPath, 0700)
	os.MkdirAll(layersPath, 0700)
	os.MkdirAll(containersPath, 0700)
}

// NewFiraaqCommand returns the root cobra.Command for Firaaq.
func NewFiraaqCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "firaaq [OPTIONS] COMMAND",
		Short:                 "A tiny tool for managing containers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PersistentPreRunE:     isRoot,
	}

	return cmd
}

// isRoot implements a cobra acceptable function and
// returns ErrNotPermitted if user is not root.
func isRoot(_ *cobra.Command, _ []string) error {
	if os.Getuid() != 0 {
		return ErrNotPermitted
	}
	return nil
}
