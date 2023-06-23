package work_packages

import "github.com/opf/openproject-cli/models"

func (dto *WorkPackageDto) convert() *models.WorkPackage {
	return &models.WorkPackage{
		Id:          dto.Id,
		Subject:     dto.Subject,
		Type:        dto.Links.Type.Title,
		Assignee:    models.Principal{Name: dto.Links.Assignee.Title},
		Status:      dto.Links.Status.Title,
		Description: dto.Description.Raw,
		LockVersion: dto.LockVersion,
	}
}

func (dto *WorkPackageCollectionDto) convert() []*models.WorkPackage {
	var workPackages = make([]*models.WorkPackage, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		workPackages[idx] = p.convert()
	}

	return workPackages
}
