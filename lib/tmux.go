package lib

import (
	"os/exec"
	"strings"
)

func TmuxHasSession(sessionName string) bool {
	cmd := exec.Command("tmux", "-L", sessionName, "ls")
	out, err := cmd.Output()
	if err != nil {
		// Handle error if command execution fails
		// For simplicity, returning false in case of error
		return false
	}

	// Split the output into lines and check if sessionName exists
	outputLines := strings.Split(string(out), "\n")
	for _, line := range outputLines {
		if strings.HasPrefix(line, sessionName+":") {
			return true
		}
	}

	return false
}
