package tools

import (
	"encoding/json"
	"fmt"
	"kubemux/lib"
	"os/exec"
	"strings"
)

// ListSessionsTool lists active tmux/zellij sessions
type ListSessionsTool struct {
	BaseTool
}

// NewListSessionsTool creates a new ListSessionsTool
func NewListSessionsTool() *ListSessionsTool {
	return &ListSessionsTool{
		BaseTool: BaseTool{
			name:        "list_sessions",
			description: "List all active kubemux/tmux/zellij sessions",
			inputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"plexer": map[string]interface{}{
						"type":        "string",
						"description": "Terminal multiplexer to use (tmux or zellij)",
						"enum":        []string{"tmux", "zellij"},
					},
				},
			},
		},
	}
}

// Execute lists all sessions
func (t *ListSessionsTool) Execute(arguments map[string]interface{}) (string, error) {
	plexer := "tmux" // default
	if p, ok := arguments["plexer"].(string); ok {
		plexer = p
	}

	var sessions []string
	var err error

	if plexer == "zellij" && lib.HasZellij() {
		sessions, err = listZellijSessions()
	} else if lib.HasTmux() {
		sessions, err = listTmuxSessions()
	} else {
		return "", fmt.Errorf("no terminal multiplexer found")
	}

	if err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"plexer":   plexer,
		"sessions": sessions,
		"count":    len(sessions),
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonResult), nil
}

func listTmuxSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return []string{}, nil // No sessions or tmux not running
	}

	sessions := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(sessions) == 1 && sessions[0] == "" {
		return []string{}, nil
	}

	return sessions, nil
}

func listZellijSessions() ([]string, error) {
	cmd := exec.Command("zellij", "list-sessions", "-n")
	output, err := cmd.Output()
	if err != nil {
		return []string{}, nil // No sessions or zellij not running
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var sessions []string
	for _, line := range lines {
		if line != "" {
			// Extract session name (first word)
			parts := strings.Fields(line)
			if len(parts) > 0 {
				sessions = append(sessions, parts[0])
			}
		}
	}

	return sessions, nil
}
