package main

import (
	"hadoop/pkg/mapper2"
	"os"
)

func main() {
	source := os.Stdin
	destination := os.Stdout
	mapper2.CustomMapper2(source, destination)
}
