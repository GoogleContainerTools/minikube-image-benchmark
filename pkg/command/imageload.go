package command

import (
	"fmt"
	"os/exec"
	"time"
)

// StartMinikubeImageLoadDocker starts minikube for docker image load.
func StartMinikubeImageLoadDocker(profile string) error {
	return startMinikube(profile)
}

// StartMinikubeImageLoadContainerd starts minikube for containerd image load.
func StartMinikubeImageLoadContainerd(profile string) error {
	return startMinikube(profile, "--container-runtime=containerd")
}

// StartMinikubeImageLoadCrio start minikube for crio image load.
func StartMinikubeImageLoadCrio(profile string) error {
	return startMinikube(profile, "--container-runtime=cri-o")
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

	// verify
	if err := verifyImage("benchmark-image", profile); err != nil {
		return 0, fmt.Errorf("image was not found after image load: %v", err)
	}

	return elapsed.Seconds(), nil
}

// ClearImageLoadCache clears out caching related to the image load method.
func ClearImageLoadCache(profile string) error {
	if err := DockerSystemPrune(); err != nil {
		return err
	}
	return minikubeDockerSystemPrune(profile)
}
