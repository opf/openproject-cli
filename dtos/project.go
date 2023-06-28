package dtos

import "github.com/opf/openproject-cli/models"

type ProjectDto struct {
	Id    int64         `json:"id"`
	Name  string        `json:"name"`
	Links *projectLinks `json:"_links"`
}

type projectLinks struct {
	Types    *LinkDto `json:"types"`
	Versions *LinkDto `json:"versions"`
}

type projectElements struct {
	Elements []*ProjectDto `json:"elements"`
}

type ProjectCollectionDto struct {
	Embedded *projectElements `json:"_embedded"`
	Type     string           `json:"_type"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *ProjectCollectionDto) Convert() []*models.Project {
	var projects = make([]*models.Project, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		projects[idx] = p.Convert()
	}

	return projects
}

func (dto *ProjectDto) Convert() *models.Project {
	return &models.Project{
		Id:   uint64(dto.Id),
		Name: dto.Name,
	}
}
