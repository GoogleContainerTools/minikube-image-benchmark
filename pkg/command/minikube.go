package command

import (
	"fmt"
	"os/exec"
)

func startMinikube(profile string, driver string) error {
	c := exec.Command("./minikube", "start", "-p", profile, "--driver", driver)
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to start minikube: %v", err)
	}

	return nil
}

func enableRegistryAddon(profile string) error {
	c := exec.Command("./minikube", "-p", profile, "addons", "enable", "registry")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to enable registry addon: %v", err)
	}

	return nil
}

// DeleteMinikube deletes the minikube cluster.
func DeleteMinikube(profile string) error {
	c := exec.Command("./minikube", "delete", "--all")
	if _, err := run(c); err != nil {
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
