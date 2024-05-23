package dtos

import "github.com/opf/openproject-cli/models"

type statusDto struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsDefault  bool   `json:"isDefault"`
	IsClosed   bool   `json:"isClosed"`
	IsReadonly bool   `json:"isReadonly"`
	Position   uint64 `json:"position"`
}

type statusElements struct {
	Elements []*statusDto `json:"elements"`
}

type StatusCollectionDto struct {
	Embedded *statusElements `json:"_embedded"`
	Type     string          `json:"_type"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *StatusCollectionDto) Convert() []*models.Status {
	var projects = make([]*models.Status, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		projects[idx] = p.Convert()
	}

	return projects
}

func (dto *statusDto) Convert() *models.Status {
	return &models.Status{
		Id:         dto.Id,
		Name:       dto.Name,
		Color:      dto.Color,
		IsDefault:  dto.IsDefault,
		IsClosed:   dto.IsClosed,
		IsReadonly: dto.IsReadonly,
		Position:   dto.Position,
	}
}
