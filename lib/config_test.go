package lib

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestERBToTemplate(t *testing.T) {
	// 定义模板字符串
	rubyTemplate := `
	name: <%= @settings["project"] %>
	root: ~/RiotGames

	socket_name: <%= @settings["project"] %>
	on_project_start:
	  - export KUBECONFIG=~/.kube/config-<%= @settings["project"] %>
	  - export TMUX_SSH_PORT="$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1])')"
	  - export TMUX_SSH_HOST="<%= @settings["host"] %>"
	startup_window: kubectl
	`

	// 定义要传递给模板的数据（这里使用 map）
	settings := map[string]string{
		"project": "YourProject",
		"host":    "YourHost",
	}
	result := RenderERB(rubyTemplate, settings)

	expectedOutput := `
	name: YourProject
	root: ~/RiotGames

	socket_name: YourProject
	on_project_start:
	  - export KUBECONFIG=~/.kube/config-YourProject
	  - export TMUX_SSH_PORT="$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1])')"
	  - export TMUX_SSH_HOST="YourHost"
	startup_window: kubectl
	`
	if result != expectedOutput {
		fmt.Printf("Unit test failed: Output does not match expected result:\n%s\n%s\n", result, expectedOutput)
	}
}

func TestParseConfig(t *testing.T) {
	yamlData := `
name: test-config
tmux: 3.2
root: /path/to/project
socket_name: test-socket
on_project_start:
  - echo "Project started"
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
`

	config, err := ParseConfig(yamlData)
	assert.NoError(t, err)

	expectedConfig := Config{
		Name:       "test-config",
		Root:       "/path/to/project",
		SocketName: "test-socket",
		OnProjectStart: []string{
			"echo \"Project started\"",
		},
		RaWWindows: []map[string]interface{}{
			{
				"editor": map[string]interface{}{
					"layout": "main-vertical",
					"panes": []interface{}{
						"vim",
						"guard",
					},
				},
			},
			{
				"server": "bundle exec rails s",
			},
			{
				"logs": "tail -f log/development.log",
			},
			{
				"proxy": map[string]interface{}{
					"layout": "main-vertical",
					"panes": []interface{}{
						"lsof -i :30001",
						"ls -alh",
						"pwd",
						"htop",
					},
				},
			},
			{
				"server": "echo $PROJ",
			},
			{
				"kubectl": "kubectl get pods",
			},
		},
		Windows: []Window{
			{
				Name:   "editor",
				Layout: "main-vertical",
				Root:   "/path/to/project",
				Panes: []Pane{
					{
						Commands: []string{"vim"},
					},
					{
						Commands: []string{"guard"},
					},
				},
			},
			{
				Name: "server",
				Root: "/path/to/project",
				Panes: []Pane{
					{
						Commands: []string{"bundle exec rails s"},
					},
				},
			},
			{
				Name: "logs",
				Root: "/path/to/project",
				Panes: []Pane{
					{
						Commands: []string{"tail -f log/development.log"},
					},
				},
			},
			{
				Name:   "proxy",
				Root:   "/path/to/project",
				Layout: "main-vertical",
				Panes: []Pane{
					{
						Commands: []string{"lsof -i :30001"},
					},
					{
						Commands: []string{"ls -alh"},
					},
					{
						Commands: []string{"pwd"},
					},
					{
						Commands: []string{"htop"},
					},
				},
			},
			{
				Name: "server",
				Root: "/path/to/project",
				Panes: []Pane{
					{
						Commands: []string{"echo $PROJ"},
					},
				},
			},
			{
				Name: "kubectl",
				Root: "/path/to/project",
				Panes: []Pane{
					{
						Commands: []string{"kubectl get pods"},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedConfig.Name, config.Name)
	assert.Equal(t, expectedConfig.Root, config.Root)
	assert.Equal(t, expectedConfig.OnProjectStart, config.OnProjectStart)
	assert.Equal(t, expectedConfig.Windows, config.Windows)
}

func TestConfigWithWindowRoot(t *testing.T) {
	yamlData := `
root: /tmp
windows:
  - first:
      layout: main-vertical
      panes:
        - vim
        - guard
  - editor:
      root: /var/run
      panes:
        - guard
`

	config, err := ParseConfig(yamlData)
	assert.NoError(t, err)

	expectedConfig := Config{
		Root: "/tmp",
		Windows: []Window{
			{
				Name:   "first",
				Layout: "main-vertical",
				Root:   "/tmp",
				Panes: []Pane{
					{
						Commands: []string{"vim"},
					},
					{
						Commands: []string{"guard"},
					},
				},
			},
			{
				Name: "editor",
				Root: "/var/run",
				Panes: []Pane{
					{
						Commands: []string{"guard"},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedConfig.Root, config.Root)
	assert.Equal(t, expectedConfig.Windows, config.Windows)
}

func TestConfigWithCommands(t *testing.T) {
	yamlData := `
root: /tmp
windows:
  - proxy:
      layout: main-vertical
      panes:
        - startup:
          - ls -alh
          - ssh -D $TMUX_SSH_PORT $TMUX_SSH_HOST
        - help:
          - pwd
  - p1:
    - ls
    - pwd
  - first:
    - vim
    - guard
  - editor:
    - guard
`

	config, err := ParseConfig(yamlData)
	assert.NoError(t, err)

	expectedConfig := Config{
		Root: "/tmp",
		Windows: []Window{
			{
				Name:   "proxy",
				Root:   "/tmp",
				Layout: "main-vertical",
				Panes: []Pane{
					{
						Commands: []string{
							"ls -alh",
							"ssh -D $TMUX_SSH_PORT $TMUX_SSH_HOST",
						},
					},
					{
						Commands: []string{
							"pwd",
						},
					},
				},
			},
			{
				Name: "p1",
				Root: "/tmp",
				Panes: []Pane{
					{
						Commands: []string{"ls"},
					},
					{
						Commands: []string{"pwd"},
					},
				},
			},
			{
				Name: "first",
				Root: "/tmp",
				Panes: []Pane{
					{
						Commands: []string{"vim"},
					},
					{
						Commands: []string{"guard"},
					},
				},
			},

			{
				Name: "editor",
				Root: "/tmp",
				Panes: []Pane{
					{
						Commands: []string{"guard"},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedConfig.Root, config.Root)
	assert.Equal(t, expectedConfig.Windows, config.Windows)
}
