package parser

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/printer"
)

func Parse[T any](body []byte) (element T) {
	err := json.Unmarshal(body, &element)
	if err != nil {
		printer.Error(err)
	}

	return
}

func IdFromLink(href string) int64 {
	split := strings.Split(href, "/")
	i, err := strconv.ParseInt(split[len(split)-1], 10, 64)
	if err != nil {
		printer.Error(err)
	}

	return i
}
