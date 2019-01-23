package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	HdfsOutputDir string
}

func ReadConfig() Config {
	configuration := Config{}
	environmentType := os.Getenv("ENV")
	fmt.Println(environmentType)
	filepath := "config." + environmentType + ".json"
	fmt.Println(filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error: File not found")
	}
	decoder := json.NewDecoder(file)
	fmt.Println(decoder)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Failed to decode configuration file")
	}
	fmt.Println(configuration)
	return configuration
}
