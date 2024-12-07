package command

import (
	"fmt"
	"kubemux/lib"
	"kubemux/lib/asset"
	kubernetes "kubemux/lib/cloud_provider"
	"kubemux/lib/cloud_provider/km_aws"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"

	"github.com/jedib0t/go-pretty/table"
	"github.com/kirsle/configdir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagProgress bool
var flagRegion string
var flagCluster string
var flagCloudProvider string

func getConfigPath(cloud string) string {
	return filepath.Join(configdir.LocalConfig("kubemux"), "kubeconfig", cloud)
}

func awsCmd(rootCmd *rootCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Use cloud provider kubernetes clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			initExtraArgs(args)
			initLogger()
			return cloudProviderExec(flagCloudProvider)
		},
	}

	cmd.Flags().StringVarP(&flagCloudProvider, "cloud", "", "", "Set the cloud provider")
	cmd.Flags().StringVarP(&flagRegion, "region", "r", "", "Set the region")
	cmd.Flags().StringVarP(&flagCluster, "cluster", "", "", "Set the cluster name")
	cmd.Flags().BoolVarP(&flagProgress, "progress", "", true, "If we show the progress bar")
	cmd.MarkFlagRequired("cloud")

	cmd.RegisterFlagCompletionFunc("cloud", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"aws"}, cobra.ShellCompDirectiveNoSpace
	})

	cmd.RegisterFlagCompletionFunc("region", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Debugging
		// go run cmd/kubemux/main.go __complete aws --region ""
		configMap, _ := traverseConfigPath(flagCloudProvider)
		regions := make([]string, 0, len(configMap))
		for region := range configMap {
			regions = append(regions, region)
		}
		return regions, cobra.ShellCompDirectiveNoSpace
	})

	cmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// go run cmd/kubemux/main.go __complete aws --region "ap-east-1" --cluster ""
		// fmt.Println("args", flagRegion, args)
		configMap, _ := traverseConfigPath(flagCloudProvider)
		if _, ok := configMap[flagRegion]; !ok {
			return nil, cobra.ShellCompDirectiveNoSpace
		}
		return configMap[flagRegion], cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

func startCluster(cloud, region, cluster string) error {
	kubeconfigPath := getKubeConfig(cloud, region, cluster)
	if _, err := os.Stat(kubeconfigPath); err != nil {
		return err
	}

	projContent := asset.KubemuxKubeconfig

	var tmuxName string
	if strings.Contains(cluster, region) {
		tmuxName = cluster
	} else {
		tmuxName = fmt.Sprintf("%s-%s", region, cluster)
	}

	projContent = lib.RenderERB(projContent, map[string]string{
		"name":       tmuxName,
		"kubeconfig": kubeconfigPath,
	})
	config, err := lib.ParseConfig(projContent)
	if err != nil {
		return err
	}
	return runTmuxWithFlags(&config)
}

func runTmuxWithFlags(config *lib.Config) error {
	config.Debug = flagDebug
	config.TmuxArgs = flagExtraArgs
	log.Debugf("Get command args: %v", flagExtraArgs)

	config.VerifyPlexerTool(flagPlexer)
	lib.RunTmux(config)
	return nil
}

func getKubeConfig(cloud, region, cluster string) string {
	kubeconfigPath := filepath.Join(
		getConfigPath(cloud),
		fmt.Sprintf("%s_kubemux_%s", region, cluster),
	)
	return kubeconfigPath
}

// Return the region and cluster name from the kubeconfig file name
func parseKubeConfigPath(kubeconfigPath string) (string, string, error) {
	// 提取文件名
	filename := filepath.Base(kubeconfigPath)
	// 去掉扩展名（如果有的话）
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	// 按下划线分割
	parts := strings.Split(filename, "_")
	if len(parts) != 3 || parts[1] != "kubemux" {
		return "", "", fmt.Errorf("invalid kubeconfig file name: %s", filename)
	}

	region := parts[0]
	cluster := parts[2]

	return region, cluster, nil
}

// It will return a map of region to cluster names
// e.g. {"us-west-2": ["cluster1", "cluster2"], "us-east-1": ["cluster3"]}
func traverseConfigPath(cloud string) (map[string][]string, error) {
	var results [][2]string
	configPath := getConfigPath(cloud)

	err := filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			region, cluster, err := parseKubeConfigPath(path)
			if err != nil {
				log.Warnf("Skipping invalid kubeconfig file: %s\n", path)
				return nil
			}
			results = append(results, [2]string{region, cluster})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	clusterMap := make(map[string][]string)
	for _, r := range results {
		if _, ok := clusterMap[r[0]]; !ok {
			clusterMap[r[0]] = []string{r[1]}
		} else {
			clusterMap[r[0]] = append(clusterMap[r[0]], r[1])
		}
	}
	return clusterMap, nil
}

func fetchAllClusters(cloud string, showProgress bool, regions []string) ([]*kubernetes.CPCluster, error) {
	// var awsProvider
	providerMap := map[string]kubernetes.CloudProvider{
		"aws": &km_aws.AWSProvider{},
	}
	awsProvider := providerMap["aws"]

	var err error
	if len(regions) == 0 {
		regions, err = awsProvider.ListRegions()
		if err != nil {
			return nil, err
		}
	}

	var setProgress func(progress int)
	if showProgress {
		bar := progressbar.Default(int64(len(regions)))
		setProgress = func(progress int) {
			bar.Set(progress)
		}
	}

	clusters, err := awsProvider.ListClusters(regions, setProgress)
	if err != nil {
		return nil, err
	}

	groupedClusters := make(map[string][]*kubernetes.CPCluster)
	for _, c := range clusters {
		awsProvider.SaveKubeconfig(c, getKubeConfig(cloud, c.Region, c.Name))
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
	return clusters, nil
}

func cloudProviderExec(cloud string) error {
	var err error
	var regions []string

	if flagRegion != "" {
		regions = []string{flagRegion}
	}

	if flagRegion != "" && flagCluster != "" {
		return startCluster(cloud, flagRegion, flagCluster)
	}
	_, err = fetchAllClusters(cloud, flagProgress, regions)

	return err
}
