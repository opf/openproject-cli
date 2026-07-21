package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type CreateOption int

const (
	CreateSubject CreateOption = iota
	CreateType
	CreateAssignee
	CreateDescription
)

var createMap = map[CreateOption]func(projectId string, workPackage *dtos.WorkPackageDto, input string) error{
	CreateSubject:     subjectCreate,
	CreateType:        typeCreate,
	CreateAssignee:    assigneeCreate,
	CreateDescription: descriptionCreate,
}

func subjectCreate(_ string, workPackage *dtos.WorkPackageDto, input string) error {
	workPackage.Subject = input

	return nil
}

func typeCreate(projectId string, workPackage *dtos.WorkPackageDto, input string) error {
	types, err := availableTypes(&dtos.LinkDto{Href: paths.Project(projectId)})
	if err != nil {
		return err
	}

	foundType := findType(input, types)
	if foundType == nil {
		printer.ErrorText("Failed to create work package type.")
		printer.Info(fmt.Sprintf(
			"No unique available type from input %s found for project %s. Please use one of the types listed below.",
			printer.Cyan(input),
			printer.Red(projectId),
		))
		printer.Types(types.Convert())
		return openerrors.ErrHandled
	}

	if workPackage.Links == nil {
		workPackage.Links = &dtos.WorkPackageLinksDto{}
	}

	workPackage.Links.Type = foundType.Links.Self

	return nil
}

func assigneeCreate(_ string, workPackage *dtos.WorkPackageDto, input string) error {
	userId, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user id %q: must be a number", input)
	}

	if workPackage.Links == nil {
		workPackage.Links = &dtos.WorkPackageLinksDto{}
	}

	workPackage.Links.Assignee = &dtos.LinkDto{Href: paths.User(userId)}
	return nil
}

func descriptionCreate(_ string, workPackage *dtos.WorkPackageDto, input string) error {
	workPackage.Description = &dtos.LongTextDto{Format: "markdown", Raw: input}
	return nil
}

func Create(projectId string, options map[CreateOption]string) (*models.WorkPackage, error) {
	return create(projectId, options)
}

func create(projectId string, options map[CreateOption]string) (*models.WorkPackage, error) {
	workPackage := dtos.WorkPackageDto{}

	for option, value := range options {
		err := createMap[option](projectId, &workPackage, value)
		if err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(workPackage)
	if err != nil {
		return nil, err
	}

	requestData := requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(data)}
	response, err := requests.Post(paths.ProjectWorkPackages(projectId), &requestData)
	if err != nil {
		return nil, err
	}

	resultingWorkPackage := parser.Parse[dtos.WorkPackageDto](response)
	return resultingWorkPackage.Convert(), nil
}
