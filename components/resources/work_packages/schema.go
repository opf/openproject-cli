package work_packages

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
)

type SchemaField struct {
	APIName  string
	Label    string
	Type     string
	Writable bool
}

type Schema struct {
	Fields []SchemaField
}

func SchemaFor(workPackage *dtos.WorkPackageDto) (*Schema, error) {
	if workPackage == nil || workPackage.Links == nil || workPackage.Links.Schema == nil {
		return &Schema{}, nil
	}

	response, err := requests.Get(workPackage.Links.Schema.Href, nil)
	if err != nil {
		return nil, err
	}

	dto := parser.Parse[dtos.WorkPackageSchemaDto](response)
	fields := make([]SchemaField, 0, len(dto.Fields))
	for apiName, field := range dto.Fields {
		fields = append(fields, SchemaField{
			APIName:  apiName,
			Label:    field.Name,
			Type:     field.Type,
			Writable: field.Writable,
		})
	}

	return &Schema{Fields: fields}, nil
}

func (schema *Schema) fieldLabels() map[string][]string {
	labels := make(map[string][]string)
	for _, field := range schema.Fields {
		labels[field.Label] = append(labels[field.Label], field.APIName)
	}
	return labels
}
