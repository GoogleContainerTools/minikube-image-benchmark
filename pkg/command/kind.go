package command

import (
	"fmt"
	"os/exec"
	"time"
)

func StartKind(profile string, args ...string) error {
	c := exec.Command("./kind", "create", "cluster")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to start kind: %v", err)
	}

	return nil
}

func RunKind(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	build := exec.Command("docker", "build", "-t", "benchmark-kind", "-f", dockerfile, ".")
	start := time.Now()
	if _, err := run(build); err != nil {
		return 0, fmt.Errorf("failed to build via kind: %v", err)
	}

	// kind load
	imageLoad := exec.Command("./kind", "load", "docker-image", "benchmark-kind:latest")
	if _, err := run(imageLoad); err != nil {
		return 0, fmt.Errorf("failed to kind load: %v", err)
	}
	elapsed := time.Now().Sub(start)

	return elapsed.Seconds(), nil
}

func ClearKindCache(profile string) error {
	return DockerSystemPrune()
}

func deleteKind() error {
	c := exec.Command("./kind", "delete", "cluster")
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to delete kind: %v", err)
	}

	return nil
}
