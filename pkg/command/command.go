// Package command contains commands that can be run to carry out the benchmarking along with branchmark setup and teardown.
package command

import (
	"fmt"
	"os/exec"
)

// run simply runs the command and returns the output, if the command fails it returns a detailed error message.
func run(cmd *exec.Cmd) (string, error) {
	o, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("\ncommand: %s\ncommand output: %s\nerr: %v", cmd.String(), string(o), err)
	}
	return string(o), nil
}

func Delete() error {
	if err := deleteMinikube(); err != nil {
		return err
	}

	return deleteKind()
}
