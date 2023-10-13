// The benchmark command is a utility that benchmarks different image build/push methods, calculates the average
// run time for each, and outputs the result to a csv file.
package main

import (
	"flag"
	"log"

	"benchmark/pkg/benchmark"
	"benchmark/pkg/command"
	"benchmark/pkg/csv"
	"benchmark/pkg/download"
)

func main() {
	runs := flag.Int("runs", 100, "number of runs per benchmark")
	profile := flag.String("profile", "benchmark", "profile to use for minikube commands")
	images := flag.String("images", "", "a comma separated list of images to benchmark")

	benchFlows := flag.String("iters", "iterative,non-iterative", "a comma separated list of flows to benchmark, options [iterative,non-iterative]")

	benchMethods := flag.String("bench-methods", "", "a comma separated list of benchmark method names")
	memory := flag.String("memory", "", "Amount of RAM to allocate to Kubernetes (format: <number>[<unit>], where unit = b, k, m or g). Use \"max\" to use the maximum amount of memory")

	flag.Parse()

	if *runs <= 0 {
		log.Fatalf("--runs must be 1 or greater")
	}

	if err := download.Files(); err != nil {
		log.Fatal(err)
	}

	defer command.Delete()

	extraMinikubeStartArgs := []string{}
	if *memory != "" {
		extraMinikubeStartArgs = append(extraMinikubeStartArgs, "--memory="+*memory)
	}
	results, err := benchmark.Run(*runs, benchmark.NewBenchMarkRunConfig(*profile, *images, *benchFlows, *benchMethods, extraMinikubeStartArgs))
	if err != nil {
		log.Printf("failed running benchmarks: %v", err)
		return
	}
	if err := csv.WriteTo(results); err != nil {
		log.Printf("failed to write to csv: %v", err)
		return
	}
}
