package main

import (
	"hadoop/pkg/mapper"
	"os"
)

func main() {
	source := os.Stdin
	destination := os.Stdout
	mapper.CustomMap(source, destination)
}
