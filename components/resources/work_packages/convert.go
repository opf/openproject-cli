package work_packages

import "github.com/opf/openproject-cli/models"

func (dto *WorkPackageDto) convert() *models.WorkPackage {
	return &models.WorkPackage{
		Id:       dto.Id,
		Subject:  dto.Subject,
		Type:     dto.Links.Type.Title,
		Assignee: dto.Links.Assignee.Title,
		Status:   dto.Links.Status.Title,
	}
}
