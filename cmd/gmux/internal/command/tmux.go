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
		Short:   "Attach gmux",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println(args)

			c.PersistentFlags().BoolVarP(&flagDeAttch, "deattch", "d", false, "If we are in debug mode")
			_, err := fmt.Printf("gmux version v%s\n", args)
			return err
		},
	}

	return cmd
}
