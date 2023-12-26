package lib

// Config 结构体用于匹配 YAML 配置文件的结构
type Config struct {
	Name    string   `yaml:"name"`
	Tmux    string   `yaml:"tmux"`
	Root    string   `yaml:"root"`
	Windows []Window `yaml:"windows"`
}
