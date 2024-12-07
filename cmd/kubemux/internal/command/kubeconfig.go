package command

import (
	"fmt"
	"kubemux/lib"
	"kubemux/lib/asset"
	"strings"

	"github.com/spf13/cobra"
)

var flagKube string

func kubeCmd(rootCmd *rootCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube",
		Short: "Display one or many resources",
		RunE:  runE,
	}

	cmd.Flags().StringVarP(&flagKube, "kube", "", "", "set the kubeconfig")
	cmd.MarkFlagRequired("kube")
	cmd.RegisterFlagCompletionFunc("kube", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return lib.GetKubeConfigList(), cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

func runE(c *cobra.Command, args []string) error {
	initExtraArgs(args)
	initLogger()
	projContent := asset.KubemuxKubeconfig
	projContent = lib.RenderERB(projContent, map[string]string{
		"name":       strings.ReplaceAll(flagKube, ".", "-"),
		"kubeconfig": fmt.Sprintf("~/.kube/%s", flagKube),
	})
	config, err := lib.ParseConfig(projContent)
	if err != nil {
		return err
	}

	return runTmuxWithFlags(&config)
}
