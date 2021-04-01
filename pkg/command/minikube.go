package command

import (
	"fmt"
	"os/exec"
)

// StartMinikube starts minikube and enables the registry addon.
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

// DeleteMinikube deletes minikube.
func DeleteMinikube(profile string) error {
	fmt.Printf("Deleting minikube...\n")
	delete := exec.Command("./minikube", "delete", "-p", profile)
	if _, err := run(delete); err != nil {
		return fmt.Errorf("failed to delete minikube: %v", err)
	}
	return nil
}

// minikube gets the IP of the running minikube instance.
func minikubeIP(profile string) (string, error) {
	c := exec.Command("./minikube", "-p", profile, "ip")
	ip, err := run(c)
	if err != nil {
		return "", fmt.Errorf("failed to get minikube ip: %v", err)
	}
	// output contains newline char, strip it out
	ip = ip[:len(ip)-1]

	return ip, nil
}
