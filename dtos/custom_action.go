package dtos

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/models"
)

type CustomActionDto struct {
	Name  string                `json:"name"`
	Links *customActionLinksDto `json:"_links"`
}

type customActionLinksDto struct {
	Self    *LinkDto `json:"self"`
	Execute *LinkDto `json:"executeImmediately"`
}

type CustomActionExecuteDto struct {
	LockVersion int              `json:"lockVersion"`
	Links       *ExecuteLinksDto `json:"_links"`
}

type ExecuteLinksDto struct {
	WorkPackage *LinkDto `json:"workPackage"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *CustomActionDto) Convert() *models.CustomAction {
	return &models.CustomAction{
		Id:   parser.IdFromLink(dto.Links.Self.Href),
		Name: dto.Name,
		}
}