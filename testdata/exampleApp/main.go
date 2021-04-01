// The main command is a utility used to generate binaries that are used in the benchmarking process.
package main

import "fmt"

// Num is a var that is meant to be set using ldflags to generate a new binary to replicate the interative build flow.
var Num = "-1"

func main() {
	fmt.Println(Num)
}
