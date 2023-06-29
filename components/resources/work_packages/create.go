package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"

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
)

var createMap = map[CreateOption]func(projectId uint64, workPackage *dtos.WorkPackageDto, input string) error{
	CreateSubject: subjectCreate,
	CreateType:    typeCreate,
}

func subjectCreate(_ uint64, workPackage *dtos.WorkPackageDto, input string) error {
	workPackage.Subject = input

	return nil
}

func typeCreate(projectId uint64, workPackage *dtos.WorkPackageDto, input string) error {
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
			printer.Red(fmt.Sprintf("#%d", projectId)),
		))

		printer.Types(types.Convert())

		return nil
	}

	if workPackage.Links == nil {
		workPackage.Links = &dtos.WorkPackageLinksDto{}
	}

	workPackage.Links.Type = foundType.Links.Self

	return nil
}

func Create(projectId uint64, options map[CreateOption]string) (*models.WorkPackage, error) {
	return create(projectId, options)
}

func create(projectId uint64, options map[CreateOption]string) (*models.WorkPackage, error) {
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
