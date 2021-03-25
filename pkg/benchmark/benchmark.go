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
var Iter = []string{ "iterative", " non-iterative"}
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
	results, err := runIterative(runs, profile, runResultsMatrix{})
	if err != nil {
		return nil, err
	}
	results, err = runNonIterative(runs, profile, results)
	if err != nil {
		return nil, err
	}
	return aggregateResults(results), nil
}

func runIterative(runs int, profile string, results runResultsMatrix) (runResultsMatrix, error) {
	for _, image := range Images {
		imageResults := map[string][]float64{}
		for _, method := range methods {
			name := method.name + Iter[0]
			fmt.Printf("\nRunning %s on %s\n", image, name)
			for i := 0; i < runs; i++ {
				if err := buildExampleApp(i); err != nil {
					return nil, err
				}
				runTime, err := method.bench(image, profile)
				if err != nil {
					return nil, fmt.Errorf("failed running benchmark %s on %s: %v", image, name, err)
				}
				imageResults[name] = append(imageResults[name], runTime)
				displayRun(i+1, runTime)
			}
			if err := method.cacheClear(profile); err != nil {
				return nil, fmt.Errorf("failed to clear cache: %v", err)
			}
		}
		results[image] = imageResults
	}
	return results, nil
}

func runNonIterative(runs int, profile string, results runResultsMatrix) (runResultsMatrix, error) {
	if err := buildExampleApp(0); err != nil {
		return nil, err
	}
	for _, image := range Images {
		imageResults := results[image]
		for _, method := range methods {
			name := method.name + Iter[1]
			fmt.Printf("\nRunning %s on %s\n", image, name)
			for i := 0; i < runs; i++ {
				runTime, err := method.bench(image, profile)
				if err != nil {
					return nil, fmt.Errorf("failed running benchmark %s on %s: %v", image, name, err)
				}
				imageResults[name] = append(imageResults[name], runTime)
				displayRun(i+1, runTime)
				if err := method.cacheClear(profile); err != nil {
					return nil, fmt.Errorf("failed to clear cache: %v", err)
				}
			}
		}
		results[image] = imageResults
	}
	return results, nil
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
