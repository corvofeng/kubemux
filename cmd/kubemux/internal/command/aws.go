package command

import (
	"fmt"
	kubernetes "kubemux/lib/cloud_provider"
	"kubemux/lib/cloud_provider/km_aws"
	"os"
	"time"

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

func awsCMDExec() {
	var awsProvider km_aws.AWSProvider
	regions, err := awsProvider.ListRegions()
	if err != nil {
		fmt.Println("Error listing regions:", err)
		return
	}

	clusters, err := awsProvider.ListClusters(regions)
	if err != nil {
		fmt.Println("Error listing clusters:", err)
		return
	}
	fmt.Println(clusters, err)
	groupedClusters := make(map[string][]*kubernetes.Cluster)
	for _, c := range clusters {
		groupedClusters[c.Region] = append(groupedClusters[c.Region], c)
	}
	bar := progressbar.Default(int64(len(clusters)))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Region", "Name", "ID", "Status"})

	for region, clusterList := range groupedClusters {
		for _, cluster := range clusterList {
			t.AppendRow(table.Row{region, cluster.Name, cluster.ID, cluster.Status})
		}
		bar.Add(1)
		time.Sleep(1000 * time.Millisecond) // Simulate loading delay
	}

	t.Render()
	// k, err := km_aws.NewEKS("eu-north-1")
	// if err != nil {
	// 	fmt.Println("Error creating EKS client:", err)
	// }

	// ch := make(chan *kubernetes.Cluster)
	// go k.GetClusters(ch)
	// clusters := []*kubernetes.Cluster{}
	// for c := range ch {
	// 	clusters = append(clusters, c)
	// 	fmt.Println(c.Name, c.Endpoint)
	// }
	// fmt.Println(clusters)

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("eu-north-1"),
	// })
	// if err != nil {
	// 	fmt.Println("Error creating session:", err)
	// 	return
	// }

	// input := &eks.ListClustersInput{}
	// svc := eks.New(sess)
	// svc.ListClustersPages(input, func(page *eks.ListClustersOutput, lastPage bool) bool {
	// 	for _, cluster := range page.Clusters {
	// 		fmt.Println(*cluster)
	// 	}
	// 	return !lastPage
	// })

	// input = &eks.DescribeClusterInput{}
	// regions, err := svc.DescribeCluster(input)
	// // regions, err := svc.DescribeClusters(&eks.DescribeClustersInput{})
	// if err != nil {
	// 	fmt.Println("Error describing EKS clusters:", err)
	// 	return
	// }

	// fmt.Println(regions)
	// 遍历每个区域
	// for regions.Next() {
	// 	// 遍历每个集群
	// 	for _, cluster := range regions.ClusterList {
	// 		fmt.Printf("Cluster: %s\n", *cluster.Name)
	// 		fmt.Printf("ARN: %s\n", *cluster.Arn)
	// 		fmt.Printf("Status: %s\n", *cluster.Status)
	// 		fmt.Println("----")
	// 	}
	// }

	// if err := regions.Err(); err != nil {
	// 	fmt.Println("Error paging through EKS clusters:", err)
	// 	return
	// }
}
