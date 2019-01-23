package main

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mitchellh/hashstructure"
	"log"
	"strconv"
	"strings"
)

var WorkerPool chan RedisWorker

type RedisWorker struct {
	client *redis.Client
}

func (worker *RedisWorker) handleRequest(payload Payload) (string){
	fmt.Println(payload)
	key, err := hashstructure.Hash(payload.Shopid + ":" + payload.Itemid, nil)
	if err != nil {
		log.Fatal("Error in hashing")
	}
	fmt.Println(key)
	reply, err := worker.client.Cmd("HGETALL", key).Map()
    fmt.Println(reply)

	if err != nil {
		log.Fatal("Error Occured")
		return "Error"
	}

	if len(reply) == 0 {
		return "Nothing found"
	}

	minDegree, _ := strconv.Atoi(payload.Degree)

	response := []Payload{}

	for _, itemAndDegree := range strings.Split(reply["itemtoDegree"], "|") {
		parts := strings.Split(itemAndDegree, "=")
		shopIdItemId, degree := parts[0], parts[1]
		degreeInt, _  := strconv.Atoi(degree)

		if degreeInt >= minDegree {
			parts = strings.Split(shopIdItemId,":")
			shopid, itemid := parts[0], parts[1]
			payloadJson, _ := json.Marshal(&Payload{Shopid:shopid, Itemid:itemid, Degree:strconv.Itoa(degreeInt)})
			fmt.Println(string(payloadJson))
			response = append(response, Payload{Shopid:shopid, Itemid:itemid, Degree:strconv.Itoa(degreeInt)})
		}

		/* Return only the top 50 items */


		if len(response) >= 50 {
			break
		}
	}
	resp, err := json.Marshal(Response{Resp: response})
	if err != nil {
		fmt.Println(err)
	}
	return string(resp)
}

func initializeWorkerPool(maxWorkers int) {
	WorkerPool = make(chan RedisWorker, maxWorkers)
	for i := 0; i < maxWorkers; i ++ {
		worker, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			log.Fatal("Error connecting client")
		}
		WorkerPool <- RedisWorker{worker}
	}
}
