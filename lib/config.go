package lib

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name       string `yaml:"name"`
	Tmux       string `yaml:"-"`
	Root       string `yaml:"root"`
	Debug      bool
	TmuxArgs   []string
	SocketName string `yaml:"socket_name"`

	/**
		There are two types of on_project_start:
		on_project_start:
	  	- export KUBECONFIG=~/.kube/config-<%= @settings["project"] %>

		on_project_start: echo "Ingame" && unset KUBECONFIG && export KUBECONFIG=~/.kube/config-bugly && echo $KUBECONFIG
	*/
	RawOnProjectStart interface{} `yaml:"on_project_start"`
	OnProjectStart    []string

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
	RawWindows []map[string]interface{} `yaml:"windows"`
	Windows    []Window                 `yaml:"-"`
}

func ParseConfig(data string) (Config, error) {
	var config Config

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		return config, err
	}

	config.Tmux = "tmux -L " + config.Name
	if config.Root == "" {
		config.Root = "~/"
	}
	if err := parseWindowConfig(&config); err != nil {
		return config, err
	}
	if err := parseProjectStart(&config); err != nil {
		return config, err
	}
	return config, nil
}

func parseWindowConfig(config *Config) error {
	config.Windows = make([]Window, len(config.RawWindows))
	for i, rawWindow := range config.RawWindows {
		window := Window{}
		for name, value := range rawWindow {
			if s, ok := value.(string); ok {
				window = Window{
					Name: name,
					Panes: []Pane{
						{Commands: []string{s}},
					}}
			} else if arr, ok := value.([]interface{}); ok {
				pane := Window{}
				for _, p := range arr {
					if cmd, ok := p.(string); ok {
						pane.Panes = append(pane.Panes, Pane{
							Commands: []string{cmd},
						})
					}
				}
				window = Window{
					Name:  name,
					Panes: pane.Panes,
				}
			} else if mp, ok := value.(map[interface{}]interface{}); ok {
				w := Window{}
				if mp["layout"] != nil {
					w.Layout = mp["layout"].(string)
				}
				if mp["root"] != nil {
					w.Root = mp["root"].(string)
				}
				if mp["panes"] != nil {
					if panes, ok := mp["panes"].([]interface{}); ok {
						for _, pane := range panes {
							paneTmp := Pane{}
							if paneCmd, ok := pane.(string); ok {
								paneTmp.Commands = append(paneTmp.Commands, paneCmd)
							} else if paneMap, ok := pane.(map[interface{}]interface{}); ok {
								for _, v := range paneMap {
									if cmds, ok := v.([]interface{}); ok {
										for _, cmd := range cmds {
											paneTmp.Commands = append(paneTmp.Commands, cmd.(string))
										}
									}
								}
							}
							w.Panes = append(w.Panes, paneTmp)
						}
					}
				}

				window = Window{
					Name:   name,
					Layout: w.Layout,
					Root:   w.Root,
					Panes:  w.Panes,
				}
			} else {
				fmt.Println("Get unknow string", value, reflect.TypeOf(value))
			}
			if window.Root == "" {
				window.Root = config.Root
			}
		}
		config.Windows[i] = window
	}
	return nil
}

func parseProjectStart(config *Config) error {
	if config.RawOnProjectStart == nil {
		return nil
	}
	arr := []string{}
	if s, ok := config.RawOnProjectStart.(string); ok {
		arr = []string{s}
	} else if cmds, ok := config.RawOnProjectStart.([]interface{}); ok {
		for _, v := range cmds {
			if s, ok := v.(string); ok {
				arr = append(arr, s)
			}
		}
	}
	config.OnProjectStart = arr
	return nil
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
