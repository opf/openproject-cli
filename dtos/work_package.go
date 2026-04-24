package dtos

import (
	"encoding/json"
	"strings"

	"github.com/opf/openproject-cli/models"
)

type WorkPackageLinksDto struct {
	Self              *LinkDto   `json:"self,omitempty"`
	AddAttachment     *LinkDto   `json:"addAttachment,omitempty"`
	Status            *LinkDto   `json:"status,omitempty"`
	Project           *LinkDto   `json:"project,omitempty"`
	Parent            *LinkDto   `json:"parent,omitempty"`
	Schema            *LinkDto   `json:"schema,omitempty"`
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
	LockVersion int                  `json:"lockVersion,omitempty"`
	CustomFields map[string]any      `json:"-"`
}

type embeddedDto struct {
	CustomActions []*CustomActionDto `json:"customActions"`
	Project       *ProjectDto        `json:"project,omitempty"`
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

func (dto *WorkPackageDto) UnmarshalJSON(data []byte) error {
	type alias WorkPackageDto
	var base alias
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	customFields := make(map[string]any)
	for key, value := range raw {
		if strings.HasPrefix(key, "customField") {
			customFields[key] = value
		}
	}

	*dto = WorkPackageDto(base)
	dto.CustomFields = customFields
	return nil
}

func (dto *WorkPackageDto) Convert() *models.WorkPackage {
	return &models.WorkPackage{
		Id:          uint64(dto.Id),
		Subject:     dto.Subject,
		Type:        linkTitle(dto.Links, func(links *WorkPackageLinksDto) *LinkDto { return links.Type }),
		Assignee:    linkTitle(dto.Links, func(links *WorkPackageLinksDto) *LinkDto { return links.Assignee }),
		Status:      linkTitle(dto.Links, func(links *WorkPackageLinksDto) *LinkDto { return links.Status }),
		Description: longTextRaw(dto.Description),
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

func linkTitle(links *WorkPackageLinksDto, selector func(*WorkPackageLinksDto) *LinkDto) string {
	if links == nil {
		return ""
	}

	link := selector(links)
	if link == nil {
		return ""
	}

	return link.Title
}

func longTextRaw(description *LongTextDto) string {
	if description == nil {
		return ""
	}

	return description.Raw
}
