package time_entries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
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

// input is validated by the caller via work_packages.ValidateIdentifier before reaching here
func workPackageCreate(entry *dtos.TimeEntryDto, input string) error {
	if entry.Links == nil {
		entry.Links = &dtos.TimeEntryLinksDto{}
	}
	entry.Links.WorkPackage = &dtos.LinkDto{Href: paths.WorkPackage(input)}
	return nil
}

func hoursCreate(entry *dtos.TimeEntryDto, input string) error {
	hours, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("invalid hours %q: must be a number", input)
	}
	if math.IsNaN(hours) || math.IsInf(hours, 0) {
		return fmt.Errorf("invalid hours %q: must be a finite number", input)
	}
	if hours <= 0 {
		return fmt.Errorf("hours must be greater than 0")
	}
	duration := hoursToISO8601(hours)
	if duration == "PT0S" {
		return fmt.Errorf("hours %q is too small: the minimum bookable duration is one second", input)
	}
	entry.Hours = duration
	return nil
}

// hoursToISO8601 formats a number of hours as an ISO 8601 time-only duration
// (e.g. PT1H30M, PT45M, PT8H). Seconds are included so sub-minute entries are
// not silently rounded away to a zero-length duration. Hours are kept as a flat
// hour count (no day/week roll-up) because OpenProject time entries are
// expressed in hours.
func hoursToISO8601(hours float64) string {
	totalSeconds := int64(math.Round(hours * 3600))
	h := totalSeconds / 3600
	m := (totalSeconds % 3600) / 60
	s := totalSeconds % 60

	out := "PT"
	if h > 0 {
		out += fmt.Sprintf("%dH", h)
	}
	if m > 0 {
		out += fmt.Sprintf("%dM", m)
	}
	if s > 0 {
		out += fmt.Sprintf("%dS", s)
	}
	if out == "PT" {
		return "PT0S"
	}
	return out
}

func activityCreate(entry *dtos.TimeEntryDto, input string) error {
	activities, err := AllActivities()
	if err != nil {
		printer.ErrorText(fmt.Sprintf("Could not fetch available activities: %s", err))
		printer.Info("Use `op time-entry list` to see existing entries and their activities.")
		return fmt.Errorf("activity lookup failed")
	}

	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("activity cannot be empty")
	}

	found, candidates := findActivity(input, activities)
	if found == nil && len(candidates) > 1 {
		printer.ErrorText(fmt.Sprintf("Activity %q is ambiguous.", input))
		printer.Info("Matching activities:")
		for _, a := range candidates {
			printer.Info(fmt.Sprintf("  - %s", a.Name))
		}
		return fmt.Errorf("activity %q is ambiguous", input)
	}
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

// findActivity resolves input to an activity: an exact (case-insensitive)
// name match wins; otherwise a prefix match is accepted only when unique.
// When the prefix is ambiguous it returns nil plus all candidates.
func findActivity(input string, activities []*models.TimeEntryActivity) (*models.TimeEntryActivity, []*models.TimeEntryActivity) {
	lower := strings.ToLower(input)
	for _, a := range activities {
		if strings.ToLower(a.Name) == lower {
			return a, nil
		}
	}

	var candidates []*models.TimeEntryActivity
	for _, a := range activities {
		if strings.HasPrefix(strings.ToLower(a.Name), lower) {
			candidates = append(candidates, a)
		}
	}
	if len(candidates) == 1 {
		return candidates[0], nil
	}
	return nil, candidates
}

func spentOnCreate(entry *dtos.TimeEntryDto, input string) error {
	if _, err := time.Parse(time.DateOnly, input); err != nil {
		return fmt.Errorf("invalid date %q: expected format YYYY-MM-DD", input)
	}
	entry.SpentOn = input
	return nil
}

func userCreate(entry *dtos.TimeEntryDto, input string) error {
	id, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
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
