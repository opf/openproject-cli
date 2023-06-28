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

	needle = 42
	if common.Contains(haystack, needle) {
		t.Errorf("Expected %v to not contain %d, but does", haystack, needle)
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

func TestFilter(t *testing.T) {
	list := []int{23, 245, 54, 132, 4325}
	filter := func(value int) bool {
		return value%2 == 0
	}

	result := common.Filter(list, filter)

	if !common.Contains(result, list[2]) || !common.Contains(result, list[3]) {
		t.Errorf("Expected %v to contain %d and %d, but does not", list, list[2], list[3])
	}
}
