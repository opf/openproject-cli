package time_entries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type CreateOption int

const (
	CreateWorkPackage CreateOption = iota
	CreateHours
	CreateActivity
	CreateSpentOn
	CreateUser
	CreateComment
)

var createMap = map[CreateOption]func(entry *dtos.TimeEntryDto, input string) error{
	CreateWorkPackage: workPackageCreate,
	CreateHours:       hoursCreate,
	CreateActivity:    activityCreate,
	CreateSpentOn:     spentOnCreate,
	CreateUser:        userCreate,
	CreateComment:     commentCreate,
}

func workPackageCreate(entry *dtos.TimeEntryDto, input string) error {
	var id uint64
	if _, err := fmt.Sscanf(input, "%d", &id); err != nil {
		return fmt.Errorf("invalid work package id %q: must be a number", input)
	}
	if entry.Links == nil {
		entry.Links = &dtos.TimeEntryLinksDto{}
	}
	entry.Links.WorkPackage = &dtos.LinkDto{Href: paths.WorkPackage(id)}
	return nil
}

func hoursCreate(entry *dtos.TimeEntryDto, input string) error {
	var hours float64
	if _, err := fmt.Sscanf(input, "%f", &hours); err != nil {
		return fmt.Errorf("invalid hours %q: must be a number", input)
	}
	if hours <= 0 {
		return fmt.Errorf("hours must be greater than 0")
	}
	entry.Hours = hoursToISO8601(hours)
	return nil
}

func hoursToISO8601(hours float64) string {
	totalMinutes := int(math.Round(hours * 60))
	h := totalMinutes / 60
	m := totalMinutes % 60
	if m == 0 {
		return fmt.Sprintf("PT%dH", h)
	}
	return fmt.Sprintf("PT%dH%dM", h, m)
}

func activityCreate(entry *dtos.TimeEntryDto, input string) error {
	activities, err := AllActivities()
	if err != nil {
		printer.ErrorText(fmt.Sprintf("Could not fetch available activities: %s", err))
		printer.Info("Use `op time-entry list` to see existing entries and their activities.")
		return fmt.Errorf("activity lookup failed")
	}

	found := findActivity(input, activities)
	if found == nil {
		printer.ErrorText(fmt.Sprintf("No activity matching %q found.", input))
		if len(activities) > 0 {
			printer.Info("Available activities:")
			for _, a := range activities {
				printer.Info(fmt.Sprintf("  - %s", a.Name))
			}
		}
		return fmt.Errorf("activity %q not found", input)
	}

	if entry.Links == nil {
		entry.Links = &dtos.TimeEntryLinksDto{}
	}
	entry.Links.Activity = &dtos.LinkDto{Href: found.Href}
	return nil
}

func findActivity(input string, activities []*models.TimeEntryActivity) *models.TimeEntryActivity {
	lower := strings.ToLower(input)
	for _, a := range activities {
		if strings.ToLower(a.Name) == lower {
			return a
		}
	}
	// partial match fallback
	for _, a := range activities {
		if strings.HasPrefix(strings.ToLower(a.Name), lower) {
			return a
		}
	}
	return nil
}

func spentOnCreate(entry *dtos.TimeEntryDto, input string) error {
	if _, err := time.Parse(time.DateOnly, input); err != nil {
		return fmt.Errorf("invalid date %q: expected format YYYY-MM-DD", input)
	}
	entry.SpentOn = input
	return nil
}

func userCreate(entry *dtos.TimeEntryDto, input string) error {
	var id uint64
	if _, err := fmt.Sscanf(input, "%d", &id); err != nil {
		return fmt.Errorf("invalid user id %q: must be a number", input)
	}
	if entry.Links == nil {
		entry.Links = &dtos.TimeEntryLinksDto{}
	}
	entry.Links.User = &dtos.LinkDto{Href: paths.User(id)}
	return nil
}

func commentCreate(entry *dtos.TimeEntryDto, input string) error {
	entry.Comment = &dtos.LongTextDto{Format: "plain", Raw: input}
	return nil
}

func AllActivities() ([]*models.TimeEntryActivity, error) {
	response, err := requests.Get(paths.TimeEntryActivities(), nil)
	if err != nil {
		return nil, err
	}
	element := parser.Parse[dtos.TimeEntryActivityCollectionDto](response)
	return element.Convert(), nil
}

func Create(options map[CreateOption]string) (*models.TimeEntry, error) {
	entry := &dtos.TimeEntryDto{}

	for option, value := range options {
		if err := createMap[option](entry, value); err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	requestData := requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(data)}
	response, err := requests.Post(paths.TimeEntries(), &requestData)
	if err != nil {
		return nil, err
	}

	result := parser.Parse[dtos.TimeEntryDto](response)
	return result.Convert(), nil
}
