package lib

type Window struct {
	Name   string
	Layout string   `yaml:"layout"`
	Root   string   `yaml:"root"`
	Panes  []string `yaml:"panes"`
}
