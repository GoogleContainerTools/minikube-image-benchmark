package command

import (
	"fmt"
	"os/exec"
	"time"
)

func StartMicrok8s(profile string, args ...string) error {
	c := exec.Command("microk8s", "start")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to start microk8s: %v", err)
	}

	return nil
}

func RunMicrok8s(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	build := exec.Command("docker", "build", "-t", "benchmark-microk8s", "-f", dockerfile, ".")
	start := time.Now()
	if _, err := run(build); err != nil {
		return 0, fmt.Errorf("failed to build via microk8s: %v", err)
	}

	// save
	args := "docker save benchmark-microk8s > benchmark-microk8s.tar"
	push := exec.Command("/bin/bash", "-c", args)
	if _, err := run(push); err != nil {
		return 0, fmt.Errorf("failed to save image via microk8s: %v", err)
	}

	// microk8s load
	imageLoad := exec.Command("microk8s", "ctr", "image", "import", "benchmark-microk8s.tar")
	if _, err := run(imageLoad); err != nil {
		return 0, fmt.Errorf("failed to microk8s load: %v", err)
	}
	elapsed := time.Now().Sub(start)

	return elapsed.Seconds(), nil
}

func ClearMicrok8sCache(profile string) error {
	return DockerSystemPrune()
}

func deleteMicrok8s() error {
	c := exec.Command("microk8s", "stop")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to stop microk8s: %v", err)
	}

	return nil
}
