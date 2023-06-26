package projects

import "github.com/opf/openproject-cli/models"

func (dto *ProjectCollectionDto) convert() []*models.Project {
	var projects = make([]*models.Project, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		projects[idx] = p.convert()
	}

	return projects
}

func (dto *ProjectDto) convert() *models.Project {
	return &models.Project{
		Id:   uint64(dto.Id),
		Name: dto.Name,
	}
}
