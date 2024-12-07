package command

import (
	"kubemux/lib"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	Logger *log.Logger
}

var flagSets []string
var flagProject string
var flagDirectory string
var flagPlexer string

// var logLevel string
var flagDebug bool

func Root(logger *log.Logger) *cobra.Command {
	rootCmd := &rootCmd{
		Logger: logger,
	}
	cmd := &cobra.Command{
		Use:   "kubemux",
		Short: "A command line tool",
		Long:  "A command line tool for handling kubemux commands",
		RunE:  rootCmd.Run,
	}

	cmd.PersistentFlags().StringSliceVar(&flagSets, "set", []string{}, "Set key-value pair")
	cmd.PersistentFlags().StringVarP(&flagProject, "project", "p", "default", "Specify the project we want to use")
	cmd.PersistentFlags().StringVarP(&flagDirectory, "directory", "", "~/.tmuxinator", "Specify the tmuxinator directory we want to use")
	cmd.PersistentFlags().StringVarP(&flagPlexer, "plexer", "", "tmux", "Specify the plexer we want to use")
	cmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "", false, "If we are in debug mode")
	cmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return lib.GetConfigList(flagDirectory), cobra.ShellCompDirectiveNoFileComp
	})

	// cmd.PersistentFlags().StringVarP(&logLevel, "lvl", "l", "INFO", "Specify log level")
	cmd.AddCommand(completionCmd(cmd))
	cmd.AddCommand(kubeCmd(rootCmd))
	// cmd.AddCommand(tmuxCmd())
	cmd.AddCommand(versionCmd())
	cmd.AddCommand(awsCmd())
	return cmd
}

func parseKeyValue(arg string) []string {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) == 2 {
		return parts
	}
	return nil
}

// inject the variables
func (c *rootCmd) ParseConfig(varMap map[string]string, configPath string) (lib.Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		c.Logger.Errorf("Error reading config file: %s %s", configPath, err)
		return lib.Config{}, err
	}
	projContent := string(content)
	projContent = lib.RenderERB(projContent, varMap)
	config, err := lib.ParseConfig(projContent)
	if err != nil {
		c.Logger.Errorf("Error parsing config file: %s %s", configPath, err)
		return lib.Config{}, err
	}

	return config, nil
}

func CreateDefaultConfig() lib.Config {
	config := lib.Config{
		Name: "kubemux-default",
		Tmux: "tmux -L kubemux-default",
		Root: "~/",
		Windows: []lib.Window{
			{
				Name: "default",
				Root: "~/",
				Panes: []lib.Pane{
					{
						Commands: []string{"ls"},
					},
				},
			},
		},
	}
	return config
}

func (c *rootCmd) Run(cmd *cobra.Command, args []string) error {
	varMap := make(map[string]string)
	for _, set := range flagSets {
		// 解析参数形式如 key=value 的形式
		keyValue := parseKeyValue(set)
		varMap[keyValue[0]] = keyValue[1]
	}
	if flagDebug {
		c.Logger.SetLevel(log.DebugLevel)
	}

	configPath := lib.ParseConfigPath(flagDirectory, flagProject)
	var config lib.Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		c.Logger.Warn("Although the kubemux works without a config file, it is recommended to create one.")
		c.Logger.Warn("Please refer to https://github.com/corvofeng/kubemux")
		config = CreateDefaultConfig()
	} else {
		if config, err = c.ParseConfig(varMap, configPath); err != nil {
			c.Logger.Errorf("Parse config error: %s", err)
			return err
		}
	}

	config.TmuxArgs = args
	config.Debug = flagDebug

	if flagPlexer == "" || flagPlexer == "tmux" {
		config.PlexerTool = lib.KTmux
	} else if flagPlexer == "zellij" {
		config.PlexerTool = lib.KZellij
	}

	lib.RunTmux(c.Logger, &config)

	return nil
}
