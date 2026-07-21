package dtos

import "github.com/opf/openproject-cli/models"

type TimeEntryActivityDto struct {
	Id      uint64                     `json:"id,omitempty"`
	Name    string                     `json:"name,omitempty"`
	Default bool                       `json:"default,omitempty"`
	Links   *timeEntryActivityLinksDto `json:"_links,omitempty"`
}

type timeEntryActivityElements struct {
	Elements []*TimeEntryActivityDto `json:"elements"`
}

type TimeEntryActivityCollectionDto struct {
	Embedded timeEntryActivityElements `json:"_embedded"`
	Type     string                    `json:"_type"`
	Total    uint64                    `json:"total"`
	Count    uint64                    `json:"count"`
}

type timeEntryActivityLinksDto struct {
	Self *LinkDto `json:"self,omitempty"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *TimeEntryActivityDto) Convert() *models.TimeEntryActivity {
	self := ""
	if dto.Links != nil && dto.Links.Self != nil {
		self = dto.Links.Self.Href
	}
	return &models.TimeEntryActivity{
		Id:      dto.Id,
		Name:    dto.Name,
		Default: dto.Default,
		Href:    self,
	}
}

func (dto *TimeEntryActivityCollectionDto) Convert() []*models.TimeEntryActivity {
	activities := make([]*models.TimeEntryActivity, len(dto.Embedded.Elements))
	for i, a := range dto.Embedded.Elements {
		activities[i] = a.Convert()
	}
	return activities
}
