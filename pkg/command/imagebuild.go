package command

import (
	"fmt"
	"os/exec"
	"time"
)

// StartMinikubeImageBuildDocker starts minikube for docker image build.
func StartMinikubeImageBuildDocker(profile string) error {
	return startMinikube(profile)
}

// StartMinikubeImageBuildContainerd starts minikube for containerd image build.
func StartMinikubeImageBuildContainerd(profile string) error {
	return startMinikube(profile, "--container-runtime=containerd")
}

// StartMinikubeImageBuildCrio start minikube for crio image build.
func StartMinikubeImageBuildCrio(profile string) error {
	return startMinikube(profile, "--container-runtime=cri-o")
}

// RunImageBuild builds the provided image using the image build method and returns the run time.
func RunImageBuild(image string, profile string) (float64, error) {
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	imageBuild := exec.Command("./minikube", "-p", profile, "image", "build", "-t", "benchmark-image-build", "-f", dockerfile, ".")
	start := time.Now()
	if _, err := run(imageBuild); err != nil {
		return 0, fmt.Errorf("failed to image build: %v", err)
	}
	elapsed := time.Now().Sub(start)

	return elapsed.Seconds(), nil
}
