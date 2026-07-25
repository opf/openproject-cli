package dtos

import (
	"encoding/json"
	"strings"
)

type WorkPackageSchemaFieldDto struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Writable bool   `json:"writable"`
}

type WorkPackageSchemaDto struct {
	Fields map[string]WorkPackageSchemaFieldDto `json:"-"`
}

func (dto *WorkPackageSchemaDto) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	fields := make(map[string]WorkPackageSchemaFieldDto)
	for key, value := range raw {
		if !strings.HasPrefix(key, "customField") {
			continue
		}

		var field WorkPackageSchemaFieldDto
		if err := json.Unmarshal(value, &field); err != nil {
			return err
		}
		fields[key] = field
	}

	dto.Fields = fields
	return nil
}
