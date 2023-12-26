package lib

import (
	"fmt"
	"testing"
)

func TestParseConfig(t *testing.T) {
	configData := []byte(`{
		"Name": "my_session",
		"Root": "~/GitRepo",
		"Windows": [
			{
				"Name": "editor",
				"Root": "~/GitRepo",
				"Panes": [
					{"Command": "vim"},
					{"Command": "ls -alh"},
					{"Command": "pwd"}
				]
			},
			{
				"Name": "shell",
				"Root": "/tmp",
				"Panes": [
					{"Command": "htop"}
				]
			}
		]
	}`)

	expectedConfig := Config{
		Name: "my_session",
		Root: "~/GitRepo",
		Windows: []Window{
			{
				Name: "editor",
				Root: "~/GitRepo",
				Panes: []Pane{
					{Command: "vim"},
					{Command: "ls -alh"},
					{Command: "pwd"},
				},
			},
			{
				Name: "shell",
				Root: "/tmp",
				Panes: []Pane{
					{Command: "htop"},
				},
			},
		},
	}

	parsedConfig, err := ParseConfig(configData)
	if err != nil {
		t.Errorf("Error parsing config: %s", err)
	}

	if fmt.Sprintf("%+v", parsedConfig) != fmt.Sprintf("%+v", expectedConfig) {
		t.Errorf("Parsed config doesn't match expected config")
	}
}

func TestParseYAMLConfig(t *testing.T) {
	yamlData := []byte(`
name: my_session
root: ~/GitRepo
windows:
  - name: editor
    root: ~/GitRepo
    panes:
      - command: vim
      - command: ls -alh
      - command: pwd
  - name: shell
    root: /tmp
    panes:
      - command: htop
  `)

	expectedConfig := Config{
		Name: "my_session",
		Root: "~/GitRepo",
		Windows: []Window{
			{
				Name: "editor",
				Root: "~/GitRepo",
				Panes: []Pane{
					{Command: "vim"},
					{Command: "ls -alh"},
					{Command: "pwd"},
				},
			},
			{
				Name: "shell",
				Root: "/tmp",
				Panes: []Pane{
					{Command: "htop"},
				},
			},
		},
	}

	parsedConfig, err := ParseYAMLConfig(yamlData)
	if err != nil {
		t.Errorf("Error parsing YAML config: %s", err)
	}

	if fmt.Sprintf("%+v", parsedConfig) != fmt.Sprintf("%+v", expectedConfig) {
		t.Errorf("Parsed YAML config doesn't match expected config:\n%s\n%s\n", parsedConfig, expectedConfig)
	}
}
