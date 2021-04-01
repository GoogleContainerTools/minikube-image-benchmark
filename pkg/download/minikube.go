package download

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// getNewestMinikube checks if a newer version of minikube exists and downloads it if there is.
func getNewestMinikube() error {
	exists, err := minikubeExists()
	if err != nil {
		return err
	}
	if !exists {
		if err := downloadFileIfNotExists("https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64", "minikube"); err != nil {
			return fmt.Errorf("failed to download minikube binary: %v", err)
		}
		if err := chmodMinikube(); err != nil {
			return err
		}
		return nil
	}
	currSHA, err := getCurrSHA()
	if err != nil {
		return err
	}
	latestSHA, err := getLatestSHA()
	if err != nil {
		return err
	}
	if currSHA == latestSHA {
		return nil
	}
	if err := os.Remove("./minikube"); err != nil {
		return fmt.Errorf("failed to delete existing minikube binary: %v", err)
	}
	fmt.Println("Newer version of minikube detected")
	return getNewestMinikube()
}

// chmodMinikube makes the minikube binary executable.
func chmodMinikube() error {
	c := exec.Command("chmod", "+x", "./minikube")
	if err := c.Run(); err != nil {
		return fmt.Errorf("failed to chmod minikube binary: %v", err)
	}
	return nil
}

// minikube exists checks if a version of minikube is already in the benchmarking directory.
func minikubeExists() (bool, error) {
	_, err := os.Stat("./minikube")
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if minikube exists: %v", err)
	}
	return true, nil
}

// getCurrSHA() gets the SHA of the current version of minikube in the benchmarking directory.
func getCurrSHA() (string, error) {
	c := exec.Command("sha256sum", "./minikube")
	o, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get minikube sha: %v", err)
	}
	currSHA := strings.Split(string(o), " ")[0]
	return currSHA, nil
}

// getLatestSHA gets the SHA of the latest version of minikube.
func getLatestSHA() (string, error) {
	c := exec.Command("/bin/bash", "-c", "curl -sL https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64.sha256")
	o, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get latest sha: %v", err)
	}
	latestSHA := string(o)
	latestSHA = latestSHA[:len(latestSHA)-1]
	return latestSHA, nil
}
