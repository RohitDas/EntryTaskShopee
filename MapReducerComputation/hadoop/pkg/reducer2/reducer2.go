package reducer2

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func CustomReducer2(source io.Reader, errors, destination io.Writer){
	input := bufio.NewScanner(source)
	var PrevKey string
	totalVal := 0
	for input.Scan(){
		line := input.Text()
		parts := strings.Split(line, "\t")
		key, val := parts[0], parts[1]

		if PrevKey == key {
			valInt, _ := strconv.Atoi(val)
			totalVal += valInt
		} else {
			if PrevKey != "" {
				fmt.Fprint(destination, fmt.Sprintf("%s\t%d\n", PrevKey, totalVal))
			}
			PrevKey = key
			totalVal, _ = strconv.Atoi(val)
		}
	}

	if PrevKey != "" {
		fmt.Fprint(destination, fmt.Sprintf("%s\t%d\n", PrevKey, totalVal))
	}
}