package main

import (
	"hadoop/pkg/reducer"
	"os"
)

func main() {
	source := os.Stdin
	destination := os.Stdout
	errors := os.Stderr
	reducer.CustomReducer(source, errors, destination)
}