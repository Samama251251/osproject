package cmd

import (
	"github.com/samama/firaaq/internal"
	"github.com/spf13/cobra"
)

// NewImagesCommand implements and returns the images command.
func NewImagesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "images",
		Short:                 "List local images",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.NoArgs,
		RunE:                  internal.Images,
	}

	return cmd
}
