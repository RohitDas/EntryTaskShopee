package main

import (
	"fmt"
	"os/exec"
	"sync"
)

/* This file is used to do stresstest for the web server */

func main() {
	var wg sync.WaitGroup
	//link := "http://localhost:8090/association/?shopid=49969001&itemid=1042094244&min_degree=1"
	link := "http://localhost:8083/association/?shopid=a&itemid=b&min_degree=1"
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func (i int) {
			defer wg.Done()
			output, _ := exec.Command("curl", link).Output()
			fmt.Printf("Response for request %d:%s \n ", i, string(output))
		}(i)

	}
	wg.Wait()
}

