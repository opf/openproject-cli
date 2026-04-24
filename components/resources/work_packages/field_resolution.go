package work_packages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/opf/openproject-cli/models"
)

var (
	ErrInvalidFieldAssignment = errors.New("invalid field assignment")
	ErrAmbiguousField         = errors.New("ambiguous field")
	ErrDuplicateField         = errors.New("duplicate field")
	ErrInvalidFieldValue      = errors.New("invalid field value")
	ErrNonWritableField       = errors.New("non-writable field")
	ErrUnknownField           = errors.New("unknown field")
	ErrUnsupportedFieldType   = errors.New("unsupported field type")
)

func resolveFieldAssignments(schema *Schema, assignments []string) (map[string]models.ResolvedField, error) {
	resolved := make(map[string]models.ResolvedField, len(assignments))
	seenAPIFields := make(map[string]struct{}, len(assignments))

	for _, assignment := range assignments {
		key, rawValue, ok := strings.Cut(assignment, "=")
		if !ok {
			return nil, fmt.Errorf("%w: %q, expected key=value", ErrInvalidFieldAssignment, assignment)
		}

		if _, exists := resolved[key]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateField, key)
		}

		field, err := resolveSchemaField(schema, key)
		if err != nil {
			return nil, err
		}

		if !field.Writable {
			return nil, fmt.Errorf("%w: %q", ErrNonWritableField, field.APIName)
		}

		if _, exists := seenAPIFields[field.APIName]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateField, field.APIName)
		}

		value, err := coerceFieldValue(field, rawValue)
		if err != nil {
			return nil, err
		}

		resolved[key] = models.ResolvedField{
			APIField: field.APIName,
			Value:    value,
		}
		seenAPIFields[field.APIName] = struct{}{}
	}

	return resolved, nil
}

func resolveSchemaField(schema *Schema, key string) (*SchemaField, error) {
	var exactAPIField *SchemaField
	var labelMatches []SchemaField

	for _, field := range schema.Fields {
		if field.APIName == key {
			fieldCopy := field
			exactAPIField = &fieldCopy
			break
		}

		if strings.EqualFold(field.Label, key) {
			labelMatches = append(labelMatches, field)
		}
	}

	if exactAPIField != nil {
		return exactAPIField, nil
	}

	switch len(labelMatches) {
	case 1:
		fieldCopy := labelMatches[0]
		return &fieldCopy, nil
	case 0:
		return nil, fmt.Errorf("%w: %q", ErrUnknownField, key)
	default:
		return nil, fmt.Errorf("%w: %q", ErrAmbiguousField, key)
	}
}

func coerceFieldValue(field *SchemaField, raw string) (any, error) {
	switch field.Type {
	case "String":
		return raw, nil
	case "Formattable":
		return map[string]any{"raw": raw}, nil
	case "Integer":
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid integer value %q for %s", ErrInvalidFieldValue, raw, field.Label)
		}
		return value, nil
	case "Float":
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid float value %q for %s", ErrInvalidFieldValue, raw, field.Label)
		}
		return value, nil
	case "Boolean":
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid boolean value %q for %s", ErrInvalidFieldValue, raw, field.Label)
		}
		return value, nil
	case "Date":
		if _, err := time.Parse("2006-01-02", raw); err != nil {
			return nil, fmt.Errorf("%w: invalid date value %q for %s", ErrInvalidFieldValue, raw, field.Label)
		}
		return raw, nil
	default:
		return nil, fmt.Errorf("%w: %s for %s", ErrUnsupportedFieldType, field.Type, field.Label)
	}
}
