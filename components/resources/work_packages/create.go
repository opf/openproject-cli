package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	componentErrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type CreateOption int

const (
	CreateSubject CreateOption = iota
	CreateParent
	CreateType
	CreateDescription
)

var createMap = map[CreateOption]func(projectId uint64, workPackage *dtos.WorkPackageDto, input string) error{
	CreateSubject:     subjectCreate,
	CreateParent:      parentCreate,
	CreateDescription: descriptionCreate,
}

func subjectCreate(_ uint64, workPackage *dtos.WorkPackageDto, input string) error {
	workPackage.Subject = input

	return nil
}

func parentCreate(_ uint64, workPackage *dtos.WorkPackageDto, input string) error {
	parentID, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return componentErrors.Custom(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", input))
	}

	if workPackage.Links == nil {
		workPackage.Links = &dtos.WorkPackageLinksDto{}
	}

	workPackage.Links.Parent = &dtos.LinkDto{Href: paths.WorkPackage(parentID)}
	return nil
}

func descriptionCreate(_ uint64, workPackage *dtos.WorkPackageDto, input string) error {
	workPackage.Description = &dtos.LongTextDto{Raw: input}

	return nil
}

func setTypeLink(workPackage *dtos.WorkPackageDto, foundType *dtos.TypeDto) {
	if workPackage.Links == nil {
		workPackage.Links = &dtos.WorkPackageLinksDto{}
	}

	workPackage.Links.Type = foundType.Links.Self
}

func Create(projectId uint64, options map[CreateOption]string) (*models.WorkPackage, error) {
	resolved, err := resolveCreate(projectId, options)
	if err != nil {
		return nil, err
	}

	return create(resolved)
}

func DryRunCreate(projectId uint64, options map[CreateOption]string) (*models.WorkPackageCreatePlan, error) {
	resolved, err := resolveCreate(projectId, options)
	if err != nil {
		return nil, err
	}

	plan := &models.WorkPackageCreatePlan{
		Valid:     true,
		Operation: "create",
		ProjectID: resolved.ProjectID,
		ParentID:  resolved.ParentID,
		WorkPackage: models.WorkPackageDraft{
			Subject:     resolved.Options[CreateSubject],
			Type:        resolved.TypeName,
			Description: resolved.Options[CreateDescription],
		},
	}

	return plan, nil
}

func create(resolved *resolvedCreate) (*models.WorkPackage, error) {
	workPackage := dtos.WorkPackageDto{}

	for _, option := range []CreateOption{CreateSubject, CreateParent, CreateDescription} {
		value, ok := resolved.Options[option]
		if !ok {
			continue
		}

		err := createMap[option](resolved.ProjectID, &workPackage, value)
		if err != nil {
			return nil, err
		}
	}

	if resolved.Type != nil {
		setTypeLink(&workPackage, resolved.Type)
	}

	data, err := json.Marshal(workPackage)
	if err != nil {
		return nil, err
	}

	requestData := requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(data)}
	response, err := requests.Post(paths.ProjectWorkPackages(resolved.ProjectID), &requestData)
	if err != nil {
		return nil, err
	}

	resultingWorkPackage := parser.Parse[dtos.WorkPackageDto](response)
	return resultingWorkPackage.Convert(), nil
}

type resolvedCreate struct {
	ProjectID uint64
	ParentID  *uint64
	Type      *dtos.TypeDto
	TypeName  string
	Options   map[CreateOption]string
}

func resolveCreate(projectId uint64, options map[CreateOption]string) (*resolvedCreate, error) {
	resolved := &resolvedCreate{
		ProjectID: projectId,
		Options:   map[CreateOption]string{},
	}

	for option, value := range options {
		resolved.Options[option] = value
	}

	if value, ok := options[CreateParent]; ok {
		parentID, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, componentErrors.Custom(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", value))
		}

		parent, err := fetch(parentID)
		if err != nil {
			return nil, err
		}

		inferredProjectID := linkID(parent.Links, func(links *dtos.WorkPackageLinksDto) *dtos.LinkDto { return links.Project })
		if inferredProjectID == nil {
			return nil, componentErrors.Custom("unable to infer project from parent work package")
		}

		if resolved.ProjectID > 0 && resolved.ProjectID != *inferredProjectID {
			return nil, componentErrors.Custom(fmt.Sprintf("conflicting project ids: --project %d does not match parent project %d", resolved.ProjectID, *inferredProjectID))
		}

		resolved.ProjectID = *inferredProjectID
		resolved.ParentID = &parentID
	}

	if resolved.ProjectID == 0 {
		return nil, componentErrors.Custom("either --project or --parent is required")
	}

	if value, ok := options[CreateType]; ok {
		foundType, err := resolveType(resolved.ProjectID, value)
		if err != nil {
			return nil, err
		}
		if foundType == nil {
			return nil, componentErrors.Custom(fmt.Sprintf("No unique available type from input %s found for project #%d.", value, resolved.ProjectID))
		}
		resolved.Type = foundType
		resolved.TypeName = foundType.Name
	}

	return resolved, nil
}

func resolveType(projectId uint64, input string) (*dtos.TypeDto, error) {
	types, err := availableTypes(&dtos.LinkDto{Href: paths.Project(projectId)})
	if err != nil {
		return nil, err
	}

	return findType(input, types), nil
}
