package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flagDeAttch bool

func tmuxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attch",
		Aliases: []string{"at"},
		Args:    cobra.ArbitraryArgs,
		Short:   "Attach kubemux",
		RunE: func(c *cobra.Command, args []string) error {
			c.PersistentFlags().BoolVarP(&flagDeAttch, "deattch", "d", false, "If we are in debug mode")
			_, err := fmt.Printf("kubemux version v%s\n", args)
			return err
		},
	}

	return cmd
}
