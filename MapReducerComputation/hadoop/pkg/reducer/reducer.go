package reducer

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func timestampGapLessThanHour(ts1 string, ts2 string) bool {
	ts1Val, _ := strconv.Atoi(ts1)
	ts2Val, _ := strconv.Atoi(ts2)

	if math.Abs(float64(ts1Val - ts2Val))/1000 <= 3600 {
		return true
	}
	return false
}

func printOut(aggregrateItems []string, destination io.Writer) {
	uidMap := make(map[string]string)
	for _, item := range aggregrateItems {
		semipart := strings.Split(item, ",")
		shopid, itemid, ts := semipart[0], semipart[1], semipart[2]
		uidMap[shopid + ":" + itemid] = ts
	}

	for k, v := range uidMap {
		for k1, v1 := range uidMap {
			if k < k1 && timestampGapLessThanHour(v,v1){
				/* Do the print here */
					fmt.Fprint(destination, fmt.Sprintf("%s\t%d\n", k + "|" + k1, 1))
				}
			}
		}
}


func CustomReducer(source io.Reader, errors, destination io.Writer){
	input := bufio.NewScanner(source)
	var PrevKey string
	aggregrateStr := []string{}
	for input.Scan(){
		line := input.Text()
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		key, val := parts[0], parts[1]

		if key == PrevKey {
			aggregrateStr = append(aggregrateStr, val)
		} else {
			if PrevKey != "" {
				printOut(aggregrateStr, destination)
			}
			PrevKey = key
			aggregrateStr = append([]string{}, val)
		}
	}
	if PrevKey != "" {
		printOut(aggregrateStr, destination)
	}

}

func WordCountReducer(source io.Reader, errors, destination io.Writer) {
	input := bufio.NewScanner(source)
	counts := make(map[string]int)

	for input.Scan() {

		line := input.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")

		word := parts[0]
		countStr := parts[1]

		count, err := strconv.Atoi(countStr)
		if err != nil {
			fmt.Fprintf(errors, "Error while processing data. Cannot convert string '%s' to number. Error: %v/n", countStr, err)
			break
		}
		counts[word] += count

	}

	for word, count := range counts {
		fmt.Fprintf(destination, fmt.Sprintf("%s\t%d\n", word, count))
	}
}
