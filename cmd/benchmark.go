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
	flag.Parse()

	if *runs <= 0 {
		log.Fatalf("--runs must be 1 or greater")
	}

	if err := download.Files(); err != nil {
		log.Fatal(err)
	}

	defer command.DeleteMinikube()

	results, err := benchmark.Run(*runs, *profile)
	if err != nil {
		log.Printf("failed running benchmarks: %v", err)
		return
	}
	if err := csv.WriteTo(results); err != nil {
		log.Printf("failed to write to csv: %v", err)
		return
	}
}
