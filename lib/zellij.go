package lib

import (
	"fmt"
	"os/exec"
)

func HasZellij() bool {
	return isCommandAvailable("zellij")
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "-V")

	if err := cmd.Run(); err != nil {
		fmt.Println("Command not found: ", name, err)
		return false
	}
	return true
}
