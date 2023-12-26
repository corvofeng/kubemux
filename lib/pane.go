package lib

type Pane struct {
	Command string
}

type MultiPane struct {
	Layout string   `yaml:"layout"`
	Panes  []string `yaml:"panes"`
}
