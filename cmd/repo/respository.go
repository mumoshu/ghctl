package repo

import (
	"io"

	"github.com/spf13/cobra"
)

// NewRepositoryCmd create new cobra command to handle repository.
func NewGetCommand(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get [resource]",
		Short:   "get resources",
		Aliases: []string{"repo"},
	}
	cmd.AddCommand(newGetRepoCmd(out, errOut))
	cmd.AddCommand(newListCmd(out, errOut))
	return cmd
}
