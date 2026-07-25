package work_packages

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func Inspect(id string) (*models.WorkPackageInspectPayload, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	schema, err := SchemaFor(workPackage)
	if err != nil {
		return nil, err
	}

	return &models.WorkPackageInspectPayload{
		WorkPackage: workPackageDetails(workPackage, schema),
		Children:    []models.WorkPackageSummary{},
	}, nil
}

func InspectWithChildren(id string) (*models.WorkPackageInspectPayload, error) {
	payload, err := Inspect(id)
	if err != nil {
		return nil, err
	}

	children, err := children(id)
	if err != nil {
		return nil, err
	}

	payload.Children = children
	return payload, nil
}

func children(parentID string) ([]models.WorkPackageSummary, error) {
	query := requests.NewUnpaginatedQuery(nil, []requests.Filter{ParentFilter(parentID)})
	response, err := requests.Get(paths.WorkPackages(), &query)
	if err != nil {
		return nil, err
	}

	collection := parser.Parse[dtos.WorkPackageCollectionDto](response)
	items := make([]models.WorkPackageSummary, 0, len(collection.Embedded.Elements))
	for _, element := range collection.Embedded.Elements {
		items = append(items, workPackageSummary(element))
	}

	return items, nil
}

func workPackageDetails(dto *dtos.WorkPackageDto, schema *Schema) models.WorkPackageDetails {
	project := models.ProjectRef{}
	if dto.Embedded != nil && dto.Embedded.Project != nil {
		project = models.ProjectRef{
			ID:         uint64(dto.Embedded.Project.Id),
			Identifier: dto.Embedded.Project.Identifier,
			Name:       dto.Embedded.Project.Name,
		}
	} else if dto.Links != nil && dto.Links.Project != nil {
		project = models.ProjectRef{
			ID:   parser.IdFromLink(dto.Links.Project.Href),
			Name: dto.Links.Project.Title,
		}
	}

	return models.WorkPackageDetails{
		ID:          uint64(dto.Id),
		Subject:     dto.Subject,
		Type:        linkTitle(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Type }),
		Status:      linkTitle(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Status }),
		Assignee:    linkTitle(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Assignee }),
		Description: longTextRaw(dto.Description),
		ParentID:    linkID(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Parent }),
		Project:     project,
		Fields:      dto.CustomFields,
		FieldLabels: schema.fieldLabels(),
	}
}

func workPackageSummary(dto *dtos.WorkPackageDto) models.WorkPackageSummary {
	return models.WorkPackageSummary{
		ID:       uint64(dto.Id),
		Subject:  dto.Subject,
		Type:     linkTitle(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Type }),
		Status:   linkTitle(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Status }),
		ParentID: linkID(dto.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Parent }),
	}
}

func linkTitle(links *dtos.WorkPackageLinksDto, selector func(*dtos.WorkPackageLinksDto) *dtos.LinkDto) string {
	if links == nil {
		return ""
	}

	link := selector(links)
	if link == nil {
		return ""
	}

	return link.Title
}

func longTextRaw(description *dtos.LongTextDto) string {
	if description == nil {
		return ""
	}

	return description.Raw
}

func linkID(links *dtos.WorkPackageLinksDto, selector func(*dtos.WorkPackageLinksDto) *dtos.LinkDto) *uint64 {
	if links == nil {
		return nil
	}

	link := selector(links)
	if link == nil || link.Href == "" {
		return nil
	}

	id := parser.IdFromLink(link.Href)
	return &id
}
