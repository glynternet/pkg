package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

// NewBashCompletion provides a completion command for the given command root
func NewBashCompletion(root *cobra.Command, w io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "generate a bash completion script",
		RunE: func(_ *cobra.Command, _ []string) error {
			return root.GenBashCompletion(w)
		},
	}
}
