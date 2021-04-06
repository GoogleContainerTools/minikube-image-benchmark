// Package benchmark runs benchmarks on different image build/push methods, and calculates the average
// run time and standard deviation for each run.
package benchmark

import (
	"fmt"
	"math"
	"os/exec"

	"benchmark/pkg/command"
)

type runResultsMatrix map[string]map[string][]float64

type aggregatedRunResult struct {
	Avg float64
	Std float64
}

// AggregatedResultsMatrix is a map containing the run results for every image method combination.
type AggregatedResultsMatrix map[string]map[string]aggregatedRunResult

type method struct {
	startMinikube func(profile string) error
	bench      func(image string, profile string) (float64, error)
	cacheClear func(profile string) error
	name       string
}

// Images is the list of all the images to use for benchmarking
var Images = []string{"alpineFewLargeFiles", "alpineFewSmallFiles", "buildpacksFewLargeFiles", "buildpacksFewSmallFiles", "openjdkFewLargeFiles", "openjdkFewSmallFiles", "ubuntuFewLargeFiles", "ubuntuFewSmallFiles"}

// Methods is a list of all the methods to use for benchmarking
var Methods = []string{"image load", "docker-env", "registry"}

// Iter contains the two flows that are benchmarked
var Iter = []string{" iterative", " non-iterative"}

var methods = []method{
	{
		command.StartMinikubeImageLoad,
		command.RunImageLoad,
		command.ClearImageLoadCache,
		"image load",
	},
	{
		command.StartMinikubeDockerEnv,
		command.RunDockerEnv,
		command.ClearDockerEnvCache,
		"docker-env",
	},
	{
		command.StartMinikubeRegistry,
		command.RunRegistry,
		command.ClearRegistryCache,
		"registry",
	}}

// Run runs all the benchmarking combinations and returns the average run time and standard deviation for each combination.
func Run(runs int, profile string) (AggregatedResultsMatrix, error) {
	modes := []func(runs int, profile string, image string, method method, imageResults map[string][]float64) error{
		runIterative,
		runNonIterative,
	}

	results := runResultsMatrix{}

	if err := buildExampleApp(0); err != nil {
		return nil, err
	}

	for _, method := range methods {
		if err := method.startMinikube(profile); err != nil {
			return nil, err
		}

		for _, mode := range modes {
			for _, image := range Images {
				imageResults := results[image]
				if imageResults == nil {
					imageResults = map[string][]float64{}
				}
				if err := mode(runs, profile, image, method, imageResults); err != nil {
					return nil, err
				}
				results[image] = imageResults
			}
		}

		if err := command.DeleteMinikube(profile); err != nil {
			return nil, err
		}
	}

	return aggregateResults(results), nil
}

// runIterative runs a benchmark using the iteratvie flow, which means changing the binary in between each run,
// mimicing an iterative flow, the cache is cleared once all the runs are complete.
func runIterative(runs int, profile string, image string, method method, imageResults map[string][]float64) error {
	name := method.name + Iter[0]
	fmt.Printf("\nRunning %s on %s\n", image, name)
	for i := 0; i < runs; i++ {
		if err := buildExampleApp(i); err != nil {
			return err
		}
		runTime, err := method.bench(image, profile)
		if err != nil {
			return fmt.Errorf("failed running benchmark %s on %s: %v", image, name, err)
		}
		imageResults[name] = append(imageResults[name], runTime)
		displayRun(i+1, runTime)
	}
	if err := method.cacheClear(profile); err != nil {
		return fmt.Errorf("failed to clear cache: %v", err)
	}

	return nil
}

// runNonIterative runs a branchmark using the non-iterative flow, which means clearing the cache after each run,
// idealy starting fresh everytime.
func runNonIterative(runs int, profile string, image string, method method, imageResults map[string][]float64) error {
	name := method.name + Iter[1]
	fmt.Printf("\nRunning %s on %s\n", image, name)
	for i := 0; i < runs; i++ {
		runTime, err := method.bench(image, profile)
		if err != nil {
			return fmt.Errorf("failed running benchmark %s on %s: %v", image, name, err)
		}
		imageResults[name] = append(imageResults[name], runTime)
		displayRun(i+1, runTime)
		if err := method.cacheClear(profile); err != nil {
			return fmt.Errorf("failed to clear cache: %v", err)
		}
	}

	return nil
}

// buildExampleApp builds the example app and sets the ldflag using the provided num.
// This allows the app the easily be changed, helping mimic the iterative workflow.
func buildExampleApp(num int) error {
	cArgs := fmt.Sprintf(`go build -o out/exampleApp -ldflags="-X 'main.Num=%d'" testdata/exampleApp/main.go`, num)
	c := exec.Command("/bin/bash", "-c", cArgs)
	o, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build example app\nerr: %v\noutput: %s", err, string(o))
	}
	return nil
}

// aggregateResults calculates the average and standard deviation from the run results
func aggregateResults(r runResultsMatrix) AggregatedResultsMatrix {
	ag := AggregatedResultsMatrix{}
	for _, image := range Images {
		imageResults := map[string]aggregatedRunResult{}
		for _, method := range Methods {
			for _, iter := range Iter {
				runs := r[image][method+iter]
				var sum, std, count float64
				for _, run := range runs {
					sum += run
					count++
				}
				avg := sum / count
				for _, run := range runs {
					std += math.Pow(run-avg, 2)
				}
				std = math.Sqrt(std / count)
				agr := aggregatedRunResult{
					Avg: avg,
					Std: std,
				}
				imageResults[method+iter] = agr
			}
		}
		ag[image] = imageResults
	}
	return ag
}

func displayRun(runNum int, runTime float64) {
	fmt.Printf("Run #%d  took %.2f seconds\n", runNum, runTime)
}
