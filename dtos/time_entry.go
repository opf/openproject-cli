package dtos

import (
	"time"

	"github.com/opf/openproject-cli/models"
	"github.com/sosodev/duration"
)

type TimeEntryDto struct {
	Id        int64              `json:"id,omitempty"`
	Comment   *LongTextDto       `json:"comment,omitempty"`
	SpentOn   string             `json:"spentOn,omitempty"`
	Hours     string             `json:"hours,omitempty"`
	Ongoing   bool               `json:"ongoing,omitempty"`
	CreatedAt string             `json:"createdAt,omitempty"`
	UpdatedAt string             `json:"updatedAt,omitempty"`
	Links     *timeEntryLinksDto `json:"_links,omitempty"`
}

type timeEntryElements struct {
	Elements []*TimeEntryDto `json:"elements"`
}

type TimeEntryCollectionDto struct {
	Embedded timeEntryElements `json:"_embedded"`
	Type     string            `json:"_type"`
	Total    uint64            `json:"total"`
	Count    uint64            `json:"count"`
}

type timeEntryLinksDto struct {
	Project     *LinkDto `json:"project"`
	WorkPackage *LinkDto `json:"workPackage,omitempty"`
	User        *LinkDto `json:"user"`
	Activity    *LinkDto `json:"activity,omitempty"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *TimeEntryDto) Convert() *models.TimeEntry {
	hours, _ := duration.Parse(dto.Hours)
	spentOn, _ := time.Parse(time.DateOnly, dto.SpentOn)

	return &models.TimeEntry{
		Id:          uint64(dto.Id),
		Comment:     dto.Comment.Raw,
		Project:     dto.Links.Project.Title,
		WorkPackage: dto.Links.WorkPackage.Title,
		SpentOn:     spentOn,
		Hours:       hours.ToTimeDuration(),
		Ongoing:     dto.Ongoing,
		User:        dto.Links.User.Title,
		Activity:    dto.Links.Activity.Title,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func (dto *TimeEntryCollectionDto) Convert() []*models.TimeEntry {
	var timeEntries = make([]*models.TimeEntry, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		timeEntries[idx] = p.Convert()
	}

	return timeEntries
}
