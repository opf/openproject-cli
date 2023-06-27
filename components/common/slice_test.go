package common_test

import (
	"github.com/opf/openproject-cli/components/common"
	"testing"
)

func TestContains(t *testing.T) {
	haystack := []int{23, 245, 54, 132, 4325}
	needle := haystack[2]

	if !common.Contains(haystack, needle) {
		t.Errorf("Expected %v to contain %d, but does not", haystack, needle)
	}
}

func TestReduce(t *testing.T) {
	list := []int{23, 245, 54, 132, 4325}
	sum := 4779

	result := common.Reduce(
		list,
		func(state int, value int) int {
			return state + value
		},
		0)

	if result != sum {
		t.Errorf("Expected %v to sum up to %d, but is actually %d", list, sum, result)
	}
}
