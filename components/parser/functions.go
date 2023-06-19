package parser

import (
	"encoding/json"
	"log"
)

func Parse[T any](body []byte) (element T) {
	err := json.Unmarshal(body, &element)

	s := string(body)
	if err != nil {
		log.Fatalf("Cannot parse content to %T. Content: %s", *new(T), s)
	}

	return
}
