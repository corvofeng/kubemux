package lib

type Window struct {
	Name  string `yaml:"name"`
	Root  string `yaml:"root"`
	Panes []Pane `yaml:"panes"`

	SocketName string `yaml:"socket_name"`
}
