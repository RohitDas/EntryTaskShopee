package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Response struct {
	Resp []Payload `json:resp`
}

type Payload struct {
	Shopid string `json:shopid`
	Itemid string `json:itemid`
	Degree string `json:degree`
}


func payloadHandler(w http.ResponseWriter, r *http.Request) {

	shopid := strings.Join(r.URL.Query()["shopid"], "")
	itemid := strings.Join(r.URL.Query()["itemid"], "")
	minDegree := strings.Join(r.URL.Query()["min_degree"], "")

	payload := Payload{Shopid: shopid, Itemid: itemid, Degree: minDegree}
	if r.Method == "POST" {
		log.Fatal("Post request  ot handled")
	} else if r.Method == "GET" {
		// Execute the command
		for {
			fmt.Println("Waiting for the workerPool")
			select {
				case worker := <- WorkerPool:
					val := worker.handleRequest(payload)
					w.Write([]byte(val))
					WorkerPool <- worker
					break
			}
			break
		}
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("A simple http app to submit jobs to the map reduce"))
}

func convToInt(str []string) (int, error) {
	return strconv.Atoi(strings.Join(str, ""))
}

func main() {
	port := os.Args[1]
	/* First initialize the worker pool */
	initializeWorkerPool(100)
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/association/", payloadHandler)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

