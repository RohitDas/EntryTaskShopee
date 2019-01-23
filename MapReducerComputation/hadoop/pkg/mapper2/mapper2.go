package mapper2

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func CustomMapper2(source io.Reader, destination io.Writer) {
	/* Read a line and create a map of item to timestamp */
	input := bufio.NewScanner(source)

	for input.Scan() {

		line := input.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		_, err := fmt.Fprint(destination, fmt.Sprintf("%s\n", line))

		if err != nil {
			continue
		}
	}
}