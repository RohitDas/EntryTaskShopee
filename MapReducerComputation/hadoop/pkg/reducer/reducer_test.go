package reducer

import (
	"bytes"
	"strings"
	"testing"
)


func TestReduceMappedElements(t *testing.T) {

	sourceTxt := "abc\t1\ndef\t1\nghi\t1\nabc\t1"
	expectedElems := []string{"abc\t2", "def\t1", "ghi\t1"}

	source := strings.NewReader(sourceTxt)
	errors := bytes.NewBuffer(nil)
	destination := bytes.NewBuffer(nil)

	Reduce(source, errors, destination)

	result := destination.String()
	for _, elem := range expectedElems {
		if !strings.Contains(result, elem) {
			t.Errorf("Result doesn't contain '%v'", elem)
		}
	}
}
