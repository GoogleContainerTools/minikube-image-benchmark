package command

import (
	"fmt"
	"os/exec"
	"time"
)

func RunImageLoad(image string, profile string) (float64, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	build := exec.Command("docker", "build", "--no-cache", "-t", "benchmark-image", "-f", dockerfile, ".")
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

func ClearImageLoadCache(profile string) error {
	// delete image from minikube to prevent caching
	deleteMinikubeImageArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker image rm benchmark-image:latest", profile)
	deleteMinikubeImage := exec.Command("/bin/bash", "-c", deleteMinikubeImageArgs)
	if _, err := run(deleteMinikubeImage); err != nil {
		return fmt.Errorf("failed to delete minikube image: %v", err)
	}

	// delete image from Docker to prevent caching
	deleteDockerImage := exec.Command("docker", "image", "rm", "benchmark-image:latest")
	if _, err := run(deleteDockerImage); err != nil {
		return fmt.Errorf("failed to delete docker image: %v", err)
	}

	// clear builder cache, must be run after the image delete
	clearBuildCache := exec.Command("docker", "builder", "prune", "-f")
	if _, err := run(clearBuildCache); err != nil {
		return fmt.Errorf("failed to clear builder cache: %v", err)
	}
	return nil
}
