package command

import (
	"fmt"
	"os/exec"
	"time"
)

func StartK3d(profile string, args ...string) error {
	c := exec.Command("k3d", "cluster", "create", "benchmark")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to start k3d: %v", err)
	}

	return nil
}

func RunK3d(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	build := exec.Command("docker", "build", "-t", "benchmark-k3d", "-f", dockerfile, ".")
	start := time.Now()
	if _, err := run(build); err != nil {
		return 0, fmt.Errorf("failed to build via k3d: %v", err)
	}

	// kind load
	imageLoad := exec.Command("k3d", "image", "import", "-c", "benchmark", "benchmark-k3d:latest")
	if _, err := run(imageLoad); err != nil {
		return 0, fmt.Errorf("failed to k3d load: %v", err)
	}
	elapsed := time.Now().Sub(start)

	return elapsed.Seconds(), nil
}

func ClearK3dCache(profile string) error {
	return DockerSystemPrune()
}

func deleteK3d() error {
	c := exec.Command("k3d", "cluster", "delete")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to delete k3d: %v", err)
	}

	return nil
}
