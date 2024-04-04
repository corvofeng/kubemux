package command

import (
	"gmux/lib"
	"os"
	"path"
	"path/filepath"
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

// var logLevel string
var flagDebug bool

func Root(logger *log.Logger) *cobra.Command {
	rootCmd := &rootCmd{
		Logger: logger,
	}
	cmd := &cobra.Command{
		Use:   "gmux",
		Short: "A command line tool",
		Long:  "A command line tool for handling gmux commands",
		RunE:  rootCmd.Run,
	}

	cmd.PersistentFlags().StringSliceVar(&flagSets, "set", []string{}, "Set key-value pair")
	cmd.PersistentFlags().StringVarP(&flagProject, "project", "p", "default", "Specify the project we want to use")
	cmd.PersistentFlags().StringVarP(&flagDirectory, "directory", "", "~/.tmuxinator", "Specify the tmuxinator directory we want to use")
	cmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "", false, "If we are in debug mode")

	// cmd.PersistentFlags().StringVarP(&logLevel, "lvl", "l", "INFO", "Specify log level")
	// cmd.AddCommand(tmuxCmd())
	// cmd.AddCommand(versionCmd())
	return cmd
}

func parseKeyValue(arg string) []string {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) == 2 {
		return parts
	}
	return nil
}

func (c *rootCmd) Run(cmd *cobra.Command, args []string) error {
	varMap := make(map[string]string)
	for _, set := range flagSets {
		// 解析参数形式如 key=value 的形式
		keyValue := parseKeyValue(set)
		varMap[keyValue[0]] = keyValue[1]
	}

	if strings.HasPrefix(flagDirectory, "~/") {
		dirname, _ := os.UserHomeDir()
		flagDirectory = filepath.Join(dirname, flagDirectory[2:])
	}
	// if we have the full path for project
	if strings.HasPrefix(flagProject, ".") || strings.HasPrefix(flagProject, "/") {
		flagDirectory = ""
	}

	if strings.HasSuffix(flagProject, "yml") {
		flagProject = strings.TrimSuffix(flagProject, ".yml")
	}

	content, err := os.ReadFile(path.Join(flagDirectory, flagProject+".yml"))
	if err != nil {
		return err
	}
	projContent := string(content)

	if len(varMap) > 0 {
		projContent = lib.RenderERB(projContent, varMap)
	}
	config, err := lib.ParseConfig(projContent)
	if err != nil {
		return err
	}
	config.Debug = flagDebug
	config.TmuxArgs = append(config.TmuxArgs, args...)
	if flagDebug {
		c.Logger.SetLevel(log.DebugLevel)
	}
	lib.RunTmux(c.Logger, &config)

	return nil
}
