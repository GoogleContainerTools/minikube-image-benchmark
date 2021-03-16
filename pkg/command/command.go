package command

import (
	"fmt"
	"os/exec"
)

func run(cmd *exec.Cmd) (string, error) {
	o, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("\ncommand: %s\ncommand output: %s\nerr: %v", cmd.String(), string(o), err)
	}
	return string(o), nil
}
