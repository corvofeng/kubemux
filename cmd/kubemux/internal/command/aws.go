package command

import (
	kubernetes "kubemux/lib/cloud_provider"
	"kubemux/lib/cloud_provider/km_aws"
	"os"

	"github.com/schollz/progressbar/v3"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagRegion string
var flagProgress bool

func awsCmd(rootCmd *rootCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Display AWS EKS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagDebug {
				log.SetLevel(log.DebugLevel)
			}
			awsCMDExec()
			return nil
		},
	}
	cmd.Flags().StringVarP(&flagRegion, "region", "r", "", "set the region")
	cmd.RegisterFlagCompletionFunc("region", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var awsProvider km_aws.AWSProvider
		regions, _ := awsProvider.ListRegions()
		return regions, cobra.ShellCompDirectiveNoSpace
	})

	cmd.Flags().BoolVarP(&flagProgress, "progress", "", true, "If we show the progress bar")
	return cmd
}

func awsCMDExec() error {
	var awsProvider km_aws.AWSProvider
	var regions []string
	var err error
	if flagRegion == "" {
		regions, err = awsProvider.ListRegions()
		if err != nil {
			return err
		}
	} else {
		regions = []string{flagRegion}
	}

	var setProgress func(progress int)
	if flagProgress {
		bar := progressbar.Default(int64(len(regions)))
		setProgress = func(progress int) {
			bar.Set(progress)
		}
	}

	clusters, err := awsProvider.ListClusters(regions, setProgress)
	if err != nil {
		return err
	}

	groupedClusters := make(map[string][]*kubernetes.Cluster)
	for _, c := range clusters {
		groupedClusters[c.Region] = append(groupedClusters[c.Region], c)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Region", "Name", "ID", "Status"})
	for region, clusterList := range groupedClusters {
		for _, cluster := range clusterList {
			t.AppendRow(table.Row{
				region, cluster.Name, cluster.ID, cluster.Status,
			})
		}
	}
	t.Render()
	return nil
}
