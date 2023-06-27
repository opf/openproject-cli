package dtos

import (
	"github.com/opf/openproject-cli/models"
)

type workPackageLinksDto struct {
	Self              *LinkDto   `json:"self"`
	AddAttachment     *LinkDto   `json:"addAttachment"`
	Status            *LinkDto   `json:"status"`
	Project           *LinkDto   `json:"project"`
	Assignee          *LinkDto   `json:"assignee"`
	Type              *LinkDto   `json:"type"`
	CustomActions     []*LinkDto `json:"customActions"`
	PrepareAttachment *LinkDto   `json:"prepareAttachment"`
}

type workPackageDescription struct {
	Raw string `json:"raw"`
}

type WorkPackageDto struct {
	Id          int64                   `json:"id,omitempty"`
	Subject     string                  `json:"subject,omitempty"`
	Links       *workPackageLinksDto    `json:"_links,omitempty"`
	Description *workPackageDescription `json:"description,omitempty"`
	Embedded    *embeddedDto            `json:"_embedded,omitempty"`
	LockVersion int                     `json:"lockVersion,omitempty"`
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
}

type CreateWorkPackageDto struct {
	Subject string `json:"subject"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *WorkPackageDto) Convert() *models.WorkPackage {
	return &models.WorkPackage{
		Id:          uint64(dto.Id),
		Subject:     dto.Subject,
		Type:        dto.Links.Type.Title,
		Assignee:    models.Principal{Name: dto.Links.Assignee.Title},
		Status:      dto.Links.Status.Title,
		Description: dto.Description.Raw,
		LockVersion: dto.LockVersion,
	}
}

func (dto *WorkPackageCollectionDto) Convert() []*models.WorkPackage {
	var workPackages = make([]*models.WorkPackage, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		workPackages[idx] = p.Convert()
	}

	return workPackages
}
