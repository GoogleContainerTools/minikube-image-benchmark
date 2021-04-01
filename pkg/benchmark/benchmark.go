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

type AggregatedResultsMatrix map[string]map[string]aggregatedRunResult

type method struct {
	bench      func(image string, profile string) (float64, error)
	cacheClear func(profile string) error
	name       string
}

var Images = []string{"alpineFewLargeFiles", "alpineFewSmallFiles", "ubuntuFewLargeFiles", "ubuntuFewSmallFiles"}
var Methods = []string{"image load", "docker-env", "registry"}
var Iter = []string{" iterative", " non-iterative"}
var methods = []method{
	{
		command.RunImageLoad,
		command.ClearImageLoadCache,
		"image load",
	},
	{
		command.RunDockerEnv,
		command.ClearDockerEnvCache,
		"docker-env",
	},
	{
		command.RunRegistry,
		command.ClearRegistryCache,
		"registry",
	}}

func Run(runs int, profile string) (AggregatedResultsMatrix, error) {
	modes := []func(runs int, profile string, image string, method method, imageResults map[string][]float64) error{
		runIterative,
		runNonIterative,
	}

	results := runResultsMatrix{}

	if err := buildExampleApp(0); err != nil {
		return nil, err
	}

	for _, mode := range modes {
		for _, image := range Images {
			imageResults := results[image]
			if imageResults == nil {
				imageResults = map[string][]float64{}
			}
			for _, method := range methods {
				if err := mode(runs, profile, image, method, imageResults); err != nil {
					return nil, err
				}
			}
			results[image] = imageResults
		}
	}

	return aggregateResults(results), nil
}

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
