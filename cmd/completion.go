package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
)

// NewBashCompletion provides a completion command for the given command root
func NewBashCompletion(root *cobra.Command, w io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: fmt.Sprintf("generate a bash completion script"),
		RunE: func(_ *cobra.Command, _ []string) error {
			return root.GenBashCompletion(w)
		},
	}
}
