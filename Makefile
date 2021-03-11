all:
	go build -o out/benchmark main.go
	GOOS=linux GOARCH=amd64 go build -o out/exampleApp testdata/exampleApp.go
