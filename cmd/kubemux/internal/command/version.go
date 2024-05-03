package command

import (
	"fmt"
	"kubemux/lib"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(c *cobra.Command, args []string) error {
			_, err := fmt.Printf("kubemux version v%s\n", lib.Version)
			return err
		},
	}

	return cmd
}
