package main

import (
	"hadoop/pkg/reducer2"
	"os"
)

func main() {
	source := os.Stdin
	destination := os.Stdout
	errors := os.Stderr
	reducer2.CustomReducer2(source, errors, destination)
}