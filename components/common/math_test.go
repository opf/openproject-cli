package common_test

import (
	"testing"

	"github.com/opf/openproject-cli/components/common"
)

func TestMax(t *testing.T) {
	input := []struct {
		x        int
		y        int
		expected int
	}{
		{x: 3, y: 3, expected: 3},
		{x: 1, y: 3, expected: 3},
		{x: 3, y: 1, expected: 3},
	}

	for _, i := range input {
		max := common.Max(i.x, i.y)
		if max != i.expected {
			t.Errorf("Expected %d, but got %d", i.expected, max)
		}
	}
}
