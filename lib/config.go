package lib

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Config 结构体用于匹配 YAML 配置文件的结构
type Config struct {
	Name    string   `yaml:"name"`
	Tmux    string   `yaml:"tmux"`
	Root    string   `yaml:"root"`
	Windows []Window `yaml:"windows"`
}

func ParseConfig(data []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	return config, err
}

func ParseYAMLConfig(data []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	return config, err
}
