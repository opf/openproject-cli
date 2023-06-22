package custom_actions

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/models"
)

func (dto *CustomActionDto) Convert() *models.CustomAction {
	return &models.CustomAction{
		Id:   parser.IdFromLink(dto.Links.Self.Href),
		Name: dto.Name,
	}
}