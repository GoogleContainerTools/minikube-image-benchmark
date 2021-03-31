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

	if err := command.StartMinikube(*profile); err != nil {
		log.Fatal(err)
	}
	defer command.DeleteMinikube(*profile)

	if err := command.SetDockerInsecureRegistry(*profile); err != nil {
		log.Printf("failed to set docker insecre registry: %v", err)
		return
	}

	// StartDockerInsecureRegistry restarts Docker, so need to start minikube again
	if err := command.StartMinikube(*profile); err != nil {
		log.Printf("failed to start minikube: %v", err)
		return
	}

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
