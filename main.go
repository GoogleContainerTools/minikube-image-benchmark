package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type runResult struct {
	runTime float64
}

type runResultsMatrix map[string]map[string][]runResult

type aggregatedRunResult struct {
	avg float64
	std float64
}

type aggregatedResultsMatrix map[string]map[string]aggregatedRunResult

type method struct {
	f    func(image string, profile string) (*runResult, error)
	name string
}

var images = []string{"alpineFewLargeFiles", "alpineFewSmallFiles", "ubuntuFewLargeFiles", "ubuntuFewSmallFiles"}
var methods = []string{"image load", "docker-env"}

func main() {
	runs := flag.Int("runs", 100, "number of runs per benchmark")
	profile := flag.String("profile", "benchmark", "profile to use for minikube commands")
	flag.Parse()

	if *runs <= 0 {
		log.Fatalf("--runs must be 1 or greater")
	}

	if err := downloadFiles(); err != nil {
		log.Fatal(err)
	}

	if err := startMinikube(*profile); err != nil {
		log.Fatal(err)
	}
	defer deleteMinikube(*profile)

	results, err := runBenchmarks(*runs, *profile)
	if err != nil {
		log.Printf("failed running benchmarks: %v", err)
		return
	}
	ag := aggregateResults(results)
	if err := writeToCSV(ag); err != nil {
		log.Printf("failed to write to csv: %v", err)
		return
	}
}

func startMinikube(profile string) error {
	fmt.Printf("Starting minikube...\n")
	start := exec.Command("./minikube", "start", "-p", "profile", "--driver", "docker")
	if _, err := runCmd(start); err != nil {
		return fmt.Errorf("failed to start minikube: %v", err)
	}
	return nil
}

func deleteMinikube(profile string) error {
	fmt.Printf("Deleting minikube...\n")
	delete := exec.Command("./minikube", "delete", "-p", profile)
	if _, err := runCmd(delete); err != nil {
		return fmt.Errorf("failed to delete minikube: %v", err)
	}
	return nil
}

func runCmd(cmd *exec.Cmd) (string, error) {
	o, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("\ncommand: %s\ncommand output: %s\nerr: %v", cmd.String(), string(o), err)
	}
	return string(o), nil
}

func runBenchmarks(runs int, profile string) (runResultsMatrix, error) {
	methods := []method{
		{
			imageLoad,
			"image load",
		},
		{
			dockerEnv,
			"docker-env",
		}}
	results := runResultsMatrix{}
	for _, image := range images {
		imageResults := map[string][]runResult{}
		for _, method := range methods {
			fmt.Printf("\nRunning %s on %s\n", image, method.name)
			for i := 0; i < runs; i++ {
				rr, err := method.f(image, profile)
				if err != nil {
					return nil, fmt.Errorf("failed running benchmark %s on %s: %v", image, method.name, err)
				}
				imageResults[method.name] = append(imageResults[method.name], *rr)
				displayRun(i+1, *rr)
			}
		}
		results[image] = imageResults
	}
	return results, nil
}

// downloadFiles downloads a 20MB & 123MB file
func downloadFiles() error {
	// 20MB file
	if err := downloadFileIfNotExists("https://golang.org/dl/go1.16.src.tar.gz", "smallFile"); err != nil {
		return err
	}
	// 123MB file
	if err := downloadFileIfNotExists("https://golang.org/dl/go1.16.linux-amd64.tar.gz", "largeFile"); err != nil {
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

// aggregateResults calculates the average and standard deviation from the run results
func aggregateResults(r runResultsMatrix) aggregatedResultsMatrix {
	ag := aggregatedResultsMatrix{}
	for _, image := range images {
		imageResults := map[string]aggregatedRunResult{}
		for _, method := range methods {
			runs := r[image][method]
			var sum, std, count float64
			for _, run := range runs {
				sum += run.runTime
				count++
			}
			avg := sum / count
			for _, run := range runs {
				std += math.Pow(run.runTime-avg, 2)
			}
			std = math.Sqrt(std / count)
			agr := aggregatedRunResult{
				avg: avg,
				std: std,
			}
			imageResults[method] = agr
		}
		ag[image] = imageResults
	}
	return ag
}

func dockerEnv(image string, profile string) (*runResult, error) {
	// build
	buildArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker build --no-cache -t benchmark-env -f testdata/Dockerfile.%s .", profile, image)
	build := exec.Command("/bin/bash", "-c", buildArgs)
	start := time.Now()
	if err := build.Run(); err != nil {
		return nil, fmt.Errorf("failed to build via docker-env: %v", err)
	}
	elapsed := time.Now().Sub(start)

	// delete image to prevent caching
	deleteArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker image rm benchmark-env:latest", profile)
	deleteImage := exec.Command("/bin/bash", "-c", deleteArgs)
	if err := deleteImage.Run(); err != nil {
		return nil, fmt.Errorf("failed to delete image: %v", err)
	}

	// clear builder cache, must be run after the image delete
	clearBuilderCacheArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker builder prune -f", profile)
	clearBuilderCache := exec.Command("/bin/bash", "-c", clearBuilderCacheArgs)
	if err := clearBuilderCache.Run(); err != nil {
		return nil, fmt.Errorf("failed to clear builder cache: %v", err)
	}
	return &runResult{runTime: elapsed.Seconds()}, nil
}

func imageLoad(image string, profile string) (*runResult, error) {
	// build
	dockerfile := fmt.Sprintf("testdata/Dockerfile.%s", image)
	build := exec.Command("docker", "build", "--no-cache", "-t", "benchmark-image", "-f", dockerfile, ".")
	start := time.Now()
	if err := build.Run(); err != nil {
		return nil, fmt.Errorf("failed to build via image load: %v", err)
	}

	// image load
	imageLoad := exec.Command("./minikube", "-p", profile, "image", "load", "benchmark-image:latest")
	if err := imageLoad.Run(); err != nil {
		return nil, fmt.Errorf("failed to image load: %v", err)
	}
	elapsed := time.Now().Sub(start)

	// verify image exists
	verifyImageArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker image ls | grep benchmark-image", profile)
	verifyImage := exec.Command("/bin/bash", "-c", verifyImageArgs)
	o, err := verifyImage.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get image list: %v", err)
	}
	if string(o) == "" {
		return nil, fmt.Errorf("image was not found after image load")
	}

	// delete image from minikube to prevent caching
	deleteMinikubeImageArgs := fmt.Sprintf("eval $(./minikube -p %s docker-env) && docker image rm benchmark-image:latest", profile)
	deleteMinikubeImage := exec.Command("/bin/bash", "-c", deleteMinikubeImageArgs)
	if err := deleteMinikubeImage.Run(); err != nil {
		return nil, fmt.Errorf("failed to delete minikube image: %v", err)
	}

	// delete image from Docker to prevent caching
	deleteDockerImage := exec.Command("docker", "image", "rm", "benchmark-image:latest")
	if err := deleteDockerImage.Run(); err != nil {
		return nil, fmt.Errorf("failed to delete docker image: %v", err)
	}

	// clear builder cache, must be run after the image delete
	clearBuildCache := exec.Command("docker", "builder", "prune", "-f")
	if err := clearBuildCache.Run(); err != nil {
		return nil, fmt.Errorf("failed to clear builder cache: %v", err)
	}
	return &runResult{runTime: elapsed.Seconds()}, nil
}

func displayRun(runNum int, rr runResult) {
	fmt.Printf("Run #%d  took %.2f seconds\n", runNum, rr.runTime)
}

func writeToCSV(ag map[string]map[string]aggregatedRunResult) error {
	records := [][]string{{"image"}}
	for _, method := range methods {
		records[0] = append(records[0], method+" average", method+" standard deviation")
	}

	f, err := os.Create("out/results.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	for _, image := range images {
		imageRecords := []string{image}
		for _, method := range methods {
			run := ag[image][method]
			imageRecords = append(imageRecords, fmt.Sprintf("%.2f", run.avg), fmt.Sprintf("%.2f", run.std))
		}
		records = append(records, imageRecords)
	}

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
}
