package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"benchmark/pkg/benchmark"
)

func WriteTo(ag benchmark.AggregatedResultsMatrix) error {
	records := [][]string{{"image"}}
	for _, method := range benchmark.Methods {
		records[0] = append(records[0], method+" average", method+" standard deviation")
	}

	f, err := os.Create("out/results.csv")
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	for _, image := range benchmark.Images {
		imageRecords := []string{image}
		for _, method := range benchmark.Methods {
			run := ag[image][method]
			imageRecords = append(imageRecords, fmt.Sprintf("%.2f", run.Avg), fmt.Sprintf("%.2f", run.Std))
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
