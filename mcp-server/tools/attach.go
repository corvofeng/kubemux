package tools

import (
	"encoding/json"
	"fmt"
)

// AttachSessionTool attaches to an existing kubemux session
type AttachSessionTool struct {
	BaseTool
}

// NewAttachSessionTool creates a new AttachSessionTool
func NewAttachSessionTool() *AttachSessionTool {
	return &AttachSessionTool{
		BaseTool: BaseTool{
			name:        "attach_session",
			description: "Attach to an existing kubemux session. Returns the command to execute.",
			inputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"session_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the session to attach to",
					},
					"plexer": map[string]interface{}{
						"type":        "string",
						"description": "Terminal multiplexer to use (tmux or zellij)",
						"enum":        []string{"tmux", "zellij"},
					},
				},
				"required": []string{"session_name"},
			},
		},
	}
}

// Execute returns the command to attach to a session
func (t *AttachSessionTool) Execute(arguments map[string]interface{}) (string, error) {
	sessionName, ok := arguments["session_name"].(string)
	if !ok {
		return "", fmt.Errorf("session_name is required")
	}

	plexer := "tmux" // default
	if p, ok := arguments["plexer"].(string); ok {
		plexer = p
	}

	var command string
	if plexer == "zellij" {
		command = fmt.Sprintf("zellij attach %s", sessionName)
	} else {
		command = fmt.Sprintf("tmux -L %s attach-session -t %s", sessionName, sessionName)
	}

	result := map[string]interface{}{
		"session_name": sessionName,
		"plexer":       plexer,
		"command":      command,
		"message":      fmt.Sprintf("To attach to the session, run: %s", command),
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonResult), nil
}

// CreateSessionTool creates a new kubemux session
type CreateSessionTool struct {
	BaseTool
}

// NewCreateSessionTool creates a new CreateSessionTool
func NewCreateSessionTool() *CreateSessionTool {
	return &CreateSessionTool{
		BaseTool: BaseTool{
			name:        "create_session",
			description: "Create a new kubemux session with specified configuration. Returns the command to execute.",
			inputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project": map[string]interface{}{
						"type":        "string",
						"description": "Project name/configuration to use (default: 'default')",
					},
					"kubeconfig": map[string]interface{}{
						"type":        "string",
						"description": "Kubeconfig file to use",
					},
					"plexer": map[string]interface{}{
						"type":        "string",
						"description": "Terminal multiplexer to use (tmux or zellij)",
						"enum":        []string{"tmux", "zellij"},
					},
					"directory": map[string]interface{}{
						"type":        "string",
						"description": "Directory for tmuxinator configurations (default: '~/.tmuxinator')",
					},
				},
			},
		},
	}
}

// Execute returns the command to create a new session
func (t *CreateSessionTool) Execute(arguments map[string]interface{}) (string, error) {
	project := "default"
	if p, ok := arguments["project"].(string); ok {
		project = p
	}

	directory := "~/.tmuxinator"
	if d, ok := arguments["directory"].(string); ok {
		directory = d
	}

	plexer := ""
	if p, ok := arguments["plexer"].(string); ok {
		plexer = p
	}

	kubeconfig := ""
	if k, ok := arguments["kubeconfig"].(string); ok {
		kubeconfig = k
	}

	// Build the kubemux command
	command := "kubemux"
	if project != "default" {
		command += fmt.Sprintf(" --project %s", project)
	}
	if directory != "~/.tmuxinator" {
		command += fmt.Sprintf(" --directory %s", directory)
	}
	if plexer != "" {
		command += fmt.Sprintf(" --plexer %s", plexer)
	}

	// If kubeconfig is specified, use the kube subcommand
	if kubeconfig != "" {
		command = fmt.Sprintf("kubemux kube --kube %s", kubeconfig)
		if plexer != "" {
			command += fmt.Sprintf(" --plexer %s", plexer)
		}
	}

	result := map[string]interface{}{
		"project":    project,
		"kubeconfig": kubeconfig,
		"plexer":     plexer,
		"command":    command,
		"message":    fmt.Sprintf("To create the session, run: %s", command),
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonResult), nil
}
