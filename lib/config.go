package lib

import (
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type TestConfig struct {
	Name string `yaml:"name"`
}
type Config struct {
	Name       string `yaml:"name"`
	Tmux       string `yaml:"-"`
	Root       string `yaml:"root"`
	SocketName string `yaml:"socket_name"`

	OnProjectStart []string `yaml:"on_project_start"`

	/**
	Since the yaml config in old tmuxinator is irregular like:
	windows:
	  - editor:
	      layout: main-vertical
	      panes:
	        - vim
	        - guard
	  - server: bundle exec rails s
	  - logs: tail -f log/development.log
	  - proxy:
	      layout: main-vertical
	      panes:
	        - lsof -i :30001
	        - ls -alh
	        - pwd
	        - htop
	  - server: echo $PROJ
	  - kubectl: kubectl get pods

	We can't parse it in the golang, so I use map to save them
	*/
	RaWWindows []map[string]interface{} `yaml:"windows"`
	Windows    []map[string]MultiPane   `yaml:"-"`
}

func ParseConfig(data string) (Config, error) {
	var config Config

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		return config, err
	}
	err := parseWindowConfig(&config)
	return config, err
}

func parseWindowConfig(config *Config) error {
	config.Windows = make([]map[string]MultiPane, len(config.RaWWindows))
	for i, rawWindow := range config.RaWWindows {
		window := make(map[string]MultiPane)
		for name, value := range rawWindow {
			if s, ok := value.(string); ok {
				window[name] = MultiPane{Panes: []string{s}}
			} else if mp, ok := value.(map[interface{}]interface{}); ok {
				pane := MultiPane{Layout: mp["layout"].(string)}
				for _, p := range mp["panes"].([]interface{}) {
					pane.Panes = append(pane.Panes, p.(string))
				}
				window[name] = MultiPane{Layout: pane.Layout, Panes: pane.Panes}
			}
		}
		config.Windows[i] = window
	}

	return nil
}

func ParseYAMLConfig(data []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	return config, err
}

type TemplateData struct {
	Settings map[string]string
}

// ConvertERBtoTemplate 将 ERB 文件内容转换为 Golang 模板字符串
func ConvertERBtoTemplate(erbContent string) string {
	// 替换 ERB 标记为模板标记
	templateStr := strings.ReplaceAll(erbContent, "<%= @settings[\"", "{{ Settings \"")
	templateStr = strings.ReplaceAll(templateStr, "\"] %>", "\" }}")

	return templateStr
}

func RenderERB(erbContent string, varMap map[string]string) string {
	goTemplateData := ConvertERBtoTemplate(erbContent)
	// 创建模板并解析模板字符串
	tmpl := template.Must(template.New("config").Funcs(template.FuncMap{
		"Settings": func(key string) string {
			return varMap[key]
		},
	}).Parse(goTemplateData))
	var result strings.Builder
	err := tmpl.Execute(&result, nil)
	if err != nil {
		panic(err)
	}
	return result.String()
}
