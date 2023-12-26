package lib

type Pane struct {
	Commands []string
}

type Window struct {
	Name   string
	Layout string `yaml:"layout"`
	Root   string `yaml:"root"`
	Panes  []Pane
}
