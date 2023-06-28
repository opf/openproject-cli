package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opf/openproject-cli/components/paths"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	workPackageUpload "github.com/opf/openproject-cli/components/resources/work_packages/upload"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type UpdateOption int

const (
	Action UpdateOption = iota
	Attach
	Subject
	Type
)

type FilterOption int

const (
	Assignee FilterOption = iota
	Version
	Project
)

func Lookup(id uint64) (*models.WorkPackage, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	return workPackage.Convert(), nil
}

func All(filterOptions *map[FilterOption]string) ([]*models.WorkPackage, error) {
	var filters []requests.Filter
	var projectId *uint64

	for updateOpt, value := range *filterOptions {
		switch updateOpt {
		case Assignee:
			filters = append(filters, AssigneeFilter(&models.Principal{Name: value}))
		case Version:
			filters = append(filters, VersionFilter(value))
		case Project:
			n, _ := strconv.ParseUint(value, 10, 64)
			projectId = &n
		}
	}

	query := requests.NewQuery(filters)

	requestUrl := paths.WorkPackages()

	if projectId != nil {
		requestUrl = paths.ProjectWorkPackages(*projectId)
	}

	response, err := requests.Get(requestUrl, &query)
	if err != nil {
		return nil, err
	}

	workPackageCollection := parser.Parse[dtos.WorkPackageCollectionDto](response)
	return workPackageCollection.Convert(), nil
}

func Create(projectId uint64, subject string) (*models.WorkPackage, error) {
	data, err := json.Marshal(dtos.CreateWorkPackageDto{Subject: subject})
	if err != nil {
		return nil, err
	}

	requestData := requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(data)}
	response, err := requests.Post(paths.ProjectWorkPackages(projectId), &requestData)
	if err != nil {
		return nil, err
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)
	return workPackage.Convert(), nil
}

func Update(id uint64, options map[UpdateOption]string) (*models.WorkPackage, error) {
	for updateOpt, value := range options {
		workPackage, err := fetch(id)
		if err != nil {
			return nil, err
		}

		switch updateOpt {
		case Action:
			err = action(workPackage, value)
		case Attach:
			err = upload(workPackage, value)
		case Subject:
			err = subject(workPackage, value)
		case Type:
			err = workPackageType(workPackage, value)
		}

		if err != nil {
			printer.Error(err)
		}
	}

	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	return workPackage.Convert(), nil
}

func fetch(id uint64) (*dtos.WorkPackageDto, error) {
	response, err := requests.Get(paths.WorkPackage(id), nil)
	if err != nil {
		return nil, err
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)
	return &workPackage, nil
}

func workPackageType(workPackage *dtos.WorkPackageDto, input string) error {
	types, err := availableTypes(workPackage)
	if err != nil {
		return err
	}

	foundType := findType(input, types)
	if foundType == nil {
		printer.ErrorText("Failed to update work package type.")
		printer.Info(fmt.Sprintf(
			"No unique available type from input %s found for project %s. Please use one of the types listed below.",
			printer.Cyan(input),
			printer.Red(fmt.Sprintf("#%d", parser.IdFromLink(workPackage.Links.Project.Href))),
		))

		printer.Types(common.Reduce(types,
			func(acc []*models.Type, dto *dtos.TypeDto) []*models.Type {
				return append(acc, dto.Convert())
			}, []*models.Type{}))
		return nil
	}

	return updateType(workPackage, foundType)
}

func updateType(workPackage *dtos.WorkPackageDto, t *dtos.TypeDto) error {
	printer.Info(fmt.Sprintf("Updating work package type to %s ...", printer.Yellow(t.Name)))

	patch := dtos.WorkPackageDto{
		LockVersion: workPackage.LockVersion,
		Links:       &dtos.WorkPackageLinksDto{Type: t.Links.Self},
	}

	marshal, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = requests.Patch(workPackage.Links.Self.Href, &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(marshal)})
	if err != nil {
		return err
	}

	printer.Done()
	return nil
}

func subject(dto *dtos.WorkPackageDto, subject string) error {
	printer.Info(fmt.Sprintf("Updating work package subject to %s ...", printer.Cyan(subject)))

	patch := dtos.WorkPackageDto{
		Subject:     subject,
		LockVersion: dto.LockVersion,
	}

	marshal, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = requests.Patch(dto.Links.Self.Href, &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(marshal)})
	if err != nil {
		return err
	}

	printer.Done()
	return nil
}

func upload(dto *dtos.WorkPackageDto, path string) error {
	if dto.Links.PrepareAttachment != nil {
		printer.ErrorText(fmt.Sprintf("Uploads to fog storages are currently not supported. :("))
	}

	printer.Info(fmt.Sprintf("Uploading %s to work package ...", printer.Yellow(filepath.Base(path))))
	link := dto.Links.AddAttachment
	reader, contentType, err := workPackageUpload.BodyReader(path)
	if err != nil {
		return err
	}

	body := &requests.RequestData{ContentType: contentType, Body: reader}
	_, err = requests.Do(link.Method, link.Href, nil, body)
	if err != nil {
		return err
	}

	printer.Done()
	return nil
}

func action(workPackage *dtos.WorkPackageDto, action string) error {
	foundAction := findAction(action, workPackage.Embedded.CustomActions)
	if foundAction == nil {
		printer.ErrorText("Failed to execute work package custom action.")
		printer.Info(fmt.Sprintf(
			"No unique available action from input %s found for work package %s. Please use one of the actions listed below.",
			printer.Cyan(action),
			printer.Red(fmt.Sprintf("#%d", workPackage.Id)),
		))
		availableActions := common.Reduce(
			workPackage.Embedded.CustomActions,
			func(acc []*models.CustomAction, dto *dtos.CustomActionDto) []*models.CustomAction {
				return append(acc, dto.Convert())
			},
			[]*models.CustomAction{},
		)
		printer.CustomActions(availableActions)
		return nil
	}

	return executeAction(workPackage, foundAction)
}

func findAction(actionInput string, availableActions []*dtos.CustomActionDto) *dtos.CustomActionDto {
	var actionAsId = false
	actionId, err := strconv.ParseUint(actionInput, 10, 64)
	if err == nil {
		actionAsId = true
	}

	var found []*dtos.CustomActionDto
	for _, act := range availableActions {
		if actionAsId && parser.IdFromLink(act.Links.Self.Href) == actionId ||
			!actionAsId && strings.ToLower(actionInput) == strings.ToLower(act.Name) {
			found = append(found, act)
		}
	}

	if len(found) == 1 {
		return found[0]
	}

	return nil
}

func executeAction(workPackage *dtos.WorkPackageDto, action *dtos.CustomActionDto) error {
	printer.Info(fmt.Sprintf("Executing action '%s' on work package [#%d] ...", action.Name, workPackage.Id))

	requestBody := &dtos.CustomActionExecuteDto{
		LockVersion: workPackage.LockVersion,
		Links:       &dtos.ExecuteLinksDto{WorkPackage: &dtos.LinkDto{Href: workPackage.Links.Self.Href}},
	}

	b, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	body := &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(b)}
	_, err = requests.Do(action.Links.Execute.Method, action.Links.Execute.Href, nil, body)
	if err != nil {
		return err
	}

	printer.Done()
	return nil
}
