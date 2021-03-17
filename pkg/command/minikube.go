package command

import (
	"fmt"
	"os/exec"
)

func StartMinikube(profile string) error {
	fmt.Printf("Starting minikube...\n")
	start := exec.Command("./minikube", "start", "-p", profile, "--driver", "docker")
	if _, err := run(start); err != nil {
		return fmt.Errorf("failed to start minikube: %v", err)
	}

	enableRegistry := exec.Command("./minikube", "-p", profile, "addons", "enable", "registry")
	if _, err := run(enableRegistry); err != nil {
		DeleteMinikube(profile)
		return fmt.Errorf("failed to enable registry addon: %v", err)
	}
	return nil
}

func DeleteMinikube(profile string) error {
	fmt.Printf("Deleting minikube...\n")
	delete := exec.Command("./minikube", "delete", "-p", profile)
	if _, err := run(delete); err != nil {
		return fmt.Errorf("failed to delete minikube: %v", err)
	}
	return nil
}
