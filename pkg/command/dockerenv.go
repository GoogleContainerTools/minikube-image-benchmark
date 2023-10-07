package command

import (
	"fmt"
	"os/exec"
	"time"
)

// StartMinikubeDockerEnv starts minikube for docker-env.
func StartMinikubeDockerEnv(profile string, args ...string) error {
	return startMinikube(profile, args...)
}

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

	// verify
	if err := verifyImage("benchmark-env", profile); err != nil {
		return 0, fmt.Errorf("image was not found after docker-env: %v", err)
	}

	return elapsed.Seconds(), nil
}
