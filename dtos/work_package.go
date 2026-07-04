package dtos

import (
	"github.com/opf/openproject-cli/models"
)

type WorkPackageLinksDto struct {
	Self              *LinkDto   `json:"self,omitempty"`
	AddAttachment     *LinkDto   `json:"addAttachment,omitempty"`
	Status            *LinkDto   `json:"status,omitempty"`
	Project           *LinkDto   `json:"project,omitempty"`
	Parent            *LinkDto   `json:"parent,omitempty"`
	Assignee          *LinkDto   `json:"assignee,omitempty"`
	Type              *LinkDto   `json:"type,omitempty"`
	CustomActions     []*LinkDto `json:"customActions,omitempty"`
	PrepareAttachment *LinkDto   `json:"prepareAttachment,omitempty"`
}

type WorkPackageDto struct {
	Id          int64                `json:"id,omitempty"`
	Subject     string               `json:"subject,omitempty"`
	Links       *WorkPackageLinksDto `json:"_links,omitempty"`
	Description *LongTextDto         `json:"description,omitempty"`
	Embedded    *embeddedDto         `json:"_embedded,omitempty"`
	LockVersion int                  `json:"lockVersion"`
}

type embeddedDto struct {
	CustomActions []*CustomActionDto `json:"customActions"`
}

type workPackageElements struct {
	Elements []*WorkPackageDto `json:"elements"`
}

type WorkPackageCollectionDto struct {
	Embedded workPackageElements `json:"_embedded"`
	Type     string              `json:"_type"`
	Total    int64               `json:"total"`
	Count    int64               `json:"count"`
	PageSize int64               `json:"pageSize"`
	Offset   int64               `json:"offset"`
}

type CreateWorkPackageDto struct {
	Subject string `json:"subject"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *WorkPackageDto) Convert() *models.WorkPackage {
	description := ""
	if dto.Description != nil {
		description = dto.Description.Raw
	}

	return &models.WorkPackage{
		Id:          uint64(dto.Id),
		Subject:     dto.Subject,
		Type:        dto.Links.Type.Title,
		Assignee:    dto.Links.Assignee.Title,
		Status:      dto.Links.Status.Title,
		Description: description,
		LockVersion: dto.LockVersion,
	}
}

func (dto *WorkPackageCollectionDto) Convert() *models.WorkPackageCollection {
	var workPackages = make([]*models.WorkPackage, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		workPackages[idx] = p.Convert()
	}

	return &models.WorkPackageCollection{
		Total:    dto.Total,
		Count:    dto.Count,
		PageSize: dto.PageSize,
		Offset:   dto.Offset,
		Items:    workPackages,
	}
}
