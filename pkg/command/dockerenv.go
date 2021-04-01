package command

import (
	"fmt"
	"os/exec"
	"time"
)

// RunDockerEnv builds the provided image using the docker-env method and returns the run time.
func RunDockerEnv(image string, profile string) (float64, error) {
	// build
	buildArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker build -t benchmark-env -f testdata/Dockerfile.%s .", profile, image)
	build := exec.Command("/bin/bash", "-c", buildArgs)
	start := time.Now()
	if _, err := run(build); err != nil {
		return 0, fmt.Errorf("failed to build via docker-env: %v", err)
	}
	elapsed := time.Now().Sub(start)

	return elapsed.Seconds(), nil
}

// ClearDockerEnvCache clears out caching related to the docker-env method.
func ClearDockerEnvCache(profile string) error {
	// delete image to prevent caching
	deleteArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker image rm benchmark-env:latest", profile)
	deleteImage := exec.Command("/bin/bash", "-c", deleteArgs)
	if _, err := run(deleteImage); err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}
	// clear builder cache, must be run after the image delete
	clearBuilderCacheArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker builder prune -f", profile)
	clearBuilderCache := exec.Command("/bin/bash", "-c", clearBuilderCacheArgs)
	if _, err := run(clearBuilderCache); err != nil {
		return fmt.Errorf("failed to clear builder cache: %v", err)
	}
	return nil
}
