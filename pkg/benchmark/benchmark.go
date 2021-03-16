package benchmark

import (
	"fmt"
	"math"

	"benchmark/pkg/command"
)

type runResultsMatrix map[string]map[string][]float64

type aggregatedRunResult struct {
	Avg float64
	Std float64
}

type AggregatedResultsMatrix map[string]map[string]aggregatedRunResult

type method struct {
	f    func(image string, profile string) (float64, error)
	name string
}

var Images = []string{"alpineFewLargeFiles", "alpineFewSmallFiles", "ubuntuFewLargeFiles", "ubuntuFewSmallFiles"}
var Methods = []string{"image load", "docker-env"}

func Run(runs int, profile string) (AggregatedResultsMatrix, error) {
	methods := []method{
		{
			command.RunImageLoad,
			"image load",
		},
		{
			command.RunDockerEnv,
			"docker-env",
		}}
	results := runResultsMatrix{}
	for _, image := range Images {
		imageResults := map[string][]float64{}
		for _, method := range methods {
			fmt.Printf("\nRunning %s on %s\n", image, method.name)
			for i := 0; i < runs; i++ {
				runTime, err := method.f(image, profile)
				if err != nil {
					return nil, fmt.Errorf("failed running benchmark %s on %s: %v", image, method.name, err)
				}
				imageResults[method.name] = append(imageResults[method.name], runTime)
				displayRun(i+1, runTime)
			}
		}
		results[image] = imageResults
	}

	return aggregateResults(results), nil
}

// aggregateResults calculates the average and standard deviation from the run results
func aggregateResults(r runResultsMatrix) AggregatedResultsMatrix {
	ag := AggregatedResultsMatrix{}
	for _, image := range Images {
		imageResults := map[string]aggregatedRunResult{}
		for _, method := range Methods {
			runs := r[image][method]
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
			imageResults[method] = agr
		}
		ag[image] = imageResults
	}
	return ag
}

func displayRun(runNum int, runTime float64) {
	fmt.Printf("Run #%d  took %.2f seconds\n", runNum, runTime)
}
