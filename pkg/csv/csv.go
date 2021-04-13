// Package csv handles writing the results of the benchmark out to a csv file.
package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"benchmark/pkg/benchmark"
)

// WriteTo writes the benchmarking results out to a csv.
func WriteTo(ag benchmark.AggregatedResultsMatrix) error {
	records := [][]string{{"image"}}
	for _, method := range benchmark.BenchMethods {
		for _, iter := range benchmark.Iter {
			records[0] = append(records[0], method.Name+iter+" average", method.Name+iter+" standard deviation")
		}
	}

	f, err := os.Create("out/results.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	for _, image := range benchmark.Images {
		imageRecords := []string{image}
		for _, method := range benchmark.BenchMethods {
			for _, iter := range benchmark.Iter {
				run := ag[image][method.Name+iter]
				imageRecords = append(imageRecords, fmt.Sprintf("%.2f", run.Avg), fmt.Sprintf("%.2f", run.Std))
			}
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
