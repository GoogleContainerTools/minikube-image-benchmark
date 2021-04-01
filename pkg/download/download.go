// Package download handles the downloading of files required for the benchmarking process.
package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Files downloads a 20MB & 123MB file
func Files() error {
	// 20MB file
	if err := downloadFileIfNotExists("https://golang.org/dl/go1.16.src.tar.gz", "smallFile"); err != nil {
		return err
	}
	// 123MB file
	if err := downloadFileIfNotExists("https://golang.org/dl/go1.16.linux-amd64.tar.gz", "largeFile"); err != nil {
		return err
	}
	if err := getNewestMinikube(); err != nil {
		return err
	}
	return nil
}

// downloadFileIfNotExists creates a file from the provided url with the provided name, if the file doesn't already exist
func downloadFileIfNotExists(url string, name string) error {
	// if file already exists skip download
	if _, err := os.Stat(name); err == nil {
		return nil
	}

	fmt.Printf("Downloading %s, please wait...\n\n", name)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file %s: %v", url, err)
	}
	defer resp.Body.Close()

	out, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", name, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to copy body to file: %v", err)
	}
	return nil
}
