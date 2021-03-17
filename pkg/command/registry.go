package command

import (
	"fmt"
	"os/exec"
	"time"
)

func RunRegistry(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	tag := fmt.Sprintf("$(./minikube -p %s ip):5000/benchmark-registry", profile)
	buildArgs := fmt.Sprintf("docker build --no-cache -t %s -f %s .", tag, dockerfile)
	build := exec.Command("/bin/bash", "-c", buildArgs)
	start := time.Now()
	if _, err := run(build); err != nil {
		return 0, fmt.Errorf("failed to build via registry: %v", err)
	}

	// push
	pushArgs := fmt.Sprintf("docker push %s", tag)
	push := exec.Command("/bin/bash", "-c", pushArgs)
	if _, err := run(push); err != nil {
		return 0, fmt.Errorf("failed to push via registry: %v", err)
	}
	elapsed := time.Now().Sub(start)

	// verify
	url := fmt.Sprintf("http://$(./minikube -p %s ip):5000/v2/_catalog", profile)
	verify := exec.Command("curl", "-X", "GET", url, "|", "grep", "benchmark-registry")
	o, err := run(verify)
	if err != nil {
		return 0, fmt.Errorf("failed to check if image was pushed successfully: %v", err)
	}
	if string(o) == "" {
		return 0, fmt.Errorf("image was not successfully pushed")
	}

	return elapsed.Seconds(), nil
}
