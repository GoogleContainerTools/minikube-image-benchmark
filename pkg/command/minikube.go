package command

import (
	"fmt"
	"os/exec"
)

func startMinikube(profile string, args ...string) error {
	a := []string{"start", "-p", profile}
	a = append(a, args...)
	c := exec.Command("./minikube", a...)
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

// deleteMinikube deletes the minikube cluster.
func deleteMinikube() error {
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

func verifyImage(image string, profile string) error {
	verifyArgs := fmt.Sprintf("./minikube -p %s image ls | grep %s", profile, image)
	verify := exec.Command("/bin/bash", "-c", verifyArgs)
	o, err := run(verify)
	if err != nil {
		return fmt.Errorf("failed to get image list: %v", err)
	}
	if string(o) == "" {
		return fmt.Errorf("image was not found")
	}

	return nil
}

// ClearDockerAndMinikubeDockerCache clears out caching related to the docker-env method.
func ClearDockerAndMinikubeDockerCache(profile string) error {
	if err := DockerSystemPrune(); err != nil {
		return err
	}
	return minikubeDockerSystemPrune(profile)
}
