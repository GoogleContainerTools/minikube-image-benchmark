all:
	go build -o out/benchmark cmd/benchmark.go
	GOOS=linux GOARCH=amd64 go build -o out/exampleApp testdata/exampleApp.go
