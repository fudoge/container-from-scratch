package main

import (
	"fmt"
	"os"
)

// dokcer run <image> <cmd> <params>
// go run main.go <image> <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])
}
