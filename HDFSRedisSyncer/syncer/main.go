package main

import (
	"bufio"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mitchellh/hashstructure"
	"httpapp/HDFSRedisSyncer/syncer/config"
	"log"
	"os"
	"os/exec"
	"strings"
)

var redisClient *redis.Client



func getListOfFiles(path string) []string{
	fmt.Println(path)
	fmt.Println("Starting to read files from path")
	cmd := "hadoop fs -ls " + path + " | sed '1d;s/  */ /g' | cut -d\\  -f8 "
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return strings.Split(string(output), "\n")
}
func readFromHDFS(path string) *bufio.Scanner {
	fmt.Println("Starting to read from hdfs")
	output, err := exec.Command("hadoop", "fs", "-cat", path).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println("Reading from hdfs complete")
	return bufio.NewScanner(strings.NewReader(string(output)))
}

func initializeClient()  {
	var err error
	redisClient, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal("Internal Error: Error connecting to the redis server")
	}

}


func updateRedisPerKey(itemA string, itemB string, degree string) {
	key, err := hashstructure.Hash(itemA, nil)
	if err != nil {
		log.Fatal(err)
	}
	reply, err := redisClient.Cmd("HGETALL", key).Map()
	if len(reply) == 0 {
		/* Key not present, Please populate */
		reply = map[string]string{
			"shopItemID": itemA,
			"itemtoDegree": itemB + "=" + degree,
		}
	} else {
		reply["itemtoDegree"] = reply["itemtoDegree"] + "|" + itemB + "=" + degree
	}
	resp := redisClient.Cmd("HMSET", key, "shopItemID", reply["shopItemID"],
		"itemtoDegree", reply["itemtoDegree"])

	if resp.Err != nil {
		log.Fatal(resp.Err)
	}
}

func updateRedis(scanner *bufio.Scanner) {
	counter := 0
	for scanner.Scan() {
		if counter % 100000 == 0 {
			fmt.Printf("Num Keys updated: %d\n",counter)
		}
		if counter >= 10000000 {
			break
		}
		counter += 1
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		itemPair, degree := parts[0], parts[1]
		items := strings.Split(itemPair, "|")
		updateRedisPerKey(items[0], items[1], degree)
		updateRedisPerKey(items[1], items[0], degree)
	}
}

func main() {
	configuration := config.ReadConfig()
	initializeClient()
	files := getListOfFiles(configuration.HdfsOutputDir)
	fmt.Println(files)
	for _, file := range files {
		fmt.Printf("Reading file: %s", file)
		scanner := readFromHDFS(file)
		updateRedis(scanner)
		fmt.Printf("Keys updated for file %s", file)
	}
}
