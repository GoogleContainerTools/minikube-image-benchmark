package command

import (
	"fmt"
	"os/exec"
)

// SetDockerInsecureRegistry sets minikube's IP in Docker's insecure registry
func SetDockerInsecureRegistry(profile string) error {
	// get minikue IP
	ip, err := minikubeIP(profile)
	if err != nil {
		return err
	}

	// create docker daemon.json
	thing := "sudo touch /etc/docker/daemon.json"
	c = exec.Command("/bin/bash", "-c", thing)
	if _, err := run(c); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	// set IP in Docker insecure registry
	args := fmt.Sprintf(`sudo tee /etc/docker/daemon.json << EOF
{
  "insecure-registries" : ["%s:5000"]
}
EOF`, ip)
	c = exec.Command("/bin/bash", "-c", args)
	if _, err = run(c); err != nil {
		return fmt.Errorf("failed to set insecure registry: %v", err)
	}

	// restart Docker so changes take effect
	c = exec.Command("sudo", "service", "docker", "restart")
	if _, err = run(c); err != nil {
		return fmt.Errorf("failed to restart docker: %v", err)
	}

	return nil
}
