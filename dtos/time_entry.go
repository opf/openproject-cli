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
	Links     *TimeEntryLinksDto `json:"_links,omitempty"`
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

type TimeEntryLinksDto struct {
	Project     *LinkDto `json:"project,omitempty"`
	WorkPackage *LinkDto `json:"workPackage,omitempty"`
	User        *LinkDto `json:"user,omitempty"`
	Activity    *LinkDto `json:"activity,omitempty"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *TimeEntryDto) Convert() *models.TimeEntry {
	hours, _ := duration.Parse(dto.Hours)
	spentOn, _ := time.Parse(time.DateOnly, dto.SpentOn)

	comment := ""
	if dto.Comment != nil {
		comment = dto.Comment.Raw
	}

	project, workPackage, user, activity := "", "", "", ""
	if dto.Links != nil {
		if dto.Links.Project != nil {
			project = dto.Links.Project.Title
		}
		if dto.Links.WorkPackage != nil {
			workPackage = dto.Links.WorkPackage.Title
		}
		if dto.Links.User != nil {
			user = dto.Links.User.Title
		}
		if dto.Links.Activity != nil {
			activity = dto.Links.Activity.Title
		}
	}

	return &models.TimeEntry{
		Id:          uint64(dto.Id),
		Comment:     comment,
		Project:     project,
		WorkPackage: workPackage,
		SpentOn:     spentOn,
		Hours:       hours.ToTimeDuration(),
		Ongoing:     dto.Ongoing,
		User:        user,
		Activity:    activity,
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
