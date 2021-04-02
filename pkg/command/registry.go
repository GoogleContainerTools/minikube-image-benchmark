package command

import (
	"fmt"
	"os/exec"
	"time"
)

// RunRegistry builds and pushes the provided image using the registry addon method and returns the run time.
func RunRegistry(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	tag := fmt.Sprintf("$(./minikube -p %s ip):5000/benchmark-registry", profile)
	buildArgs := fmt.Sprintf("docker build -t %s -f %s .", tag, dockerfile)
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
	ip, err := minikubeIP(profile)
	if err != nil {
		return 0, err
	}

	verifyArgs := fmt.Sprintf("curl http://%s:5000/v2/_catalog | grep benchmark-registry", ip)
	verify := exec.Command("/bin/bash", "-c", verifyArgs)
	o, err := run(verify)
	if err != nil {
		return 0, fmt.Errorf("failed to check if image was pushed successfully: %v", err)
	}
	if string(o) == "" {
		return 0, fmt.Errorf("image was not successfully pushed")
	}

	return elapsed.Seconds(), nil
}

// ClearRegistryCache clears out caching related to the registry addon method.
func ClearRegistryCache(profile string) error {
	return DockerSystemPrune()
}
