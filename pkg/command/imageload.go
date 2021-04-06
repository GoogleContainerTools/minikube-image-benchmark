package command

import (
	"fmt"
	"os/exec"
	"time"
)

// StartMinikubeImageLoad starts minikube for docker-env.
func StartMinikubeImageLoad(profile string) error {
	return startMinikube(profile, "docker")
}

// RunImageLoad builds the provided image using the image load method and returns the run time.
func RunImageLoad(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	build := exec.Command("docker", "build", "-t", "benchmark-image", "-f", dockerfile, ".")
	start := time.Now()
	if _, err := run(build); err != nil {
		return 0, fmt.Errorf("failed to build via image load: %v", err)
	}

	// image load
	imageLoad := exec.Command("./minikube", "-p", profile, "image", "load", "benchmark-image:latest")
	if _, err := run(imageLoad); err != nil {
		return 0, fmt.Errorf("failed to image load: %v", err)
	}
	elapsed := time.Now().Sub(start)

	// verify image exists
	verifyImageArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker image ls | grep benchmark-image", profile)
	verifyImage := exec.Command("/bin/bash", "-c", verifyImageArgs)
	o, err := run(verifyImage)
	if err != nil {
		return 0, fmt.Errorf("failed to get image list: %v", err)
	}
	if string(o) == "" {
		return 0, fmt.Errorf("image was not found after image load")
	}

	return elapsed.Seconds(), nil
}

// ClearImageLoadCache clears out caching related to the image load method.
func ClearImageLoadCache(profile string) error {
	if err := minikubeDockerSystemPrune(profile); err != nil {
		return err
	}
	return DockerSystemPrune()
}
