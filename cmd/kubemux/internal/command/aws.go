package command

import (
	"fmt"
	kubernetes "kubemux/lib/cloud_provider"
	"kubemux/lib/cloud_provider/km_aws"
	"os"

	"github.com/schollz/progressbar/v3"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func awsCmd(rootCmd *rootCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Display one or many resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagDebug {
				log.SetLevel(log.DebugLevel)
			}
			awsCMDExec()
			return nil
		},
	}
	return cmd
}

func awsCMDExec() error {
	var awsProvider km_aws.AWSProvider
	regions, err := awsProvider.ListRegions()
	if err != nil {
		fmt.Println("Error listing regions:", err)
		return err
	}

	bar := progressbar.Default(int64(len(regions)))
	clusters, err := awsProvider.ListClusters(regions, func(progress int) {
		bar.Set(progress)
	})
	if err != nil {
		fmt.Println("Error listing clusters:", err)
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
			t.AppendRow(table.Row{region, cluster.Name, cluster.ID, cluster.Status})
		}
	}
	t.Render()
	return nil
}
