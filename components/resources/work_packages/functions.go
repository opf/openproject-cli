package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	workPackageTypes "github.com/opf/openproject-cli/components/resources/work_packages/types"
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

const apiPath = "api/v3"
const workPackagesPath = apiPath + "/work_packages"

func Lookup(id int64) *models.WorkPackage {
	return fetch(id).Convert()
}

func All(principal *models.Principal) []*models.WorkPackage {
	var filters []requests.Filter

	if principal != nil {
		filters = append(filters, AssigneeFilter(principal))
	}

	query := requests.NewQuery(filters)

	status, response := requests.Get(workPackagesPath, &query)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	workPackageCollection := parser.Parse[dtos.WorkPackageCollectionDto](response)
	return workPackageCollection.Convert()
}

func Create(projectId uint64, subject string) *models.WorkPackage {
	data, err := json.Marshal(dtos.CreateWorkPackageDto{Subject: subject})
	if err != nil {
		printer.Error(err)
	}

	requestData := requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(data)}

	status, response := requests.Post(
		filepath.Join(apiPath, "projects", strconv.FormatUint(projectId, 10), "dtos"),
		&requestData,
	)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)

	return workPackage.Convert()
}

func Update(id int64, opts map[UpdateOption]string) {
	workPackage := fetch(id)

	for updateOpt, value := range opts {
		switch updateOpt {
		case Action:
			action(workPackage, value)
		case Attach:
			upload(workPackage, value)
		case Subject:
			subject(workPackage, value)
		case Type:
			workPackageType(workPackage, value)
		}
	}
}

func fetch(id int64) *dtos.WorkPackageDto {
	status, response := requests.Get(filepath.Join(workPackagesPath, strconv.FormatInt(id, 10)), nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)
	return &workPackage
}

func workPackageType(workPackage *dtos.WorkPackageDto, input string) {
	availableTypes := workPackageTypes.AvailableTypes(workPackage)

	status, response := requests.Get(workPackage.Links.Project.Href, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	foundType := findType(input, availableTypes)
	if foundType == nil {
		printer.Info(fmt.Sprintf(
			"No unique available type from input '%s' found for project [#%d]. Please use one of the types listed below.",
			input,
			parser.IdFromLink(workPackage.Links.Project.Href),
		))
		//		availableActions := common.Reduce(
		//			workPackage.Embedded.CustomActions,
		//			func(acc []*models.CustomAction, dto *dtos.CustomActionDto) []*models.CustomAction {
		//				return append(acc, dto.Convert())
		//			},
		//			[]*models.CustomAction{},
		//		)
		//		printer.CustomActions(availableActions)
		return
	}

	printer.ErrorText("not yet implemented")
}

func findType(input string, availableTypes []*models.Type) *models.Type {
	return nil
}

func subject(dto *dtos.WorkPackageDto, subject string) {
	patch := dtos.WorkPackageDto{
		Subject:     subject,
		LockVersion: dto.LockVersion,
	}

	marshal, err := json.Marshal(patch)
	if err != nil {
		printer.Error(err)
	}

	status, response := requests.Patch(dto.Links.Self.Href, &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(marshal)})
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	} else {
		printer.Done()
	}
}

func upload(dto *dtos.WorkPackageDto, path string) {
	if dto.Links.PrepareAttachment != nil {
		printer.ErrorText(fmt.Sprintf("Uploads to fog storages are currently not supported. :("))
	}

	link := dto.Links.AddAttachment
	reader, contentType, err := workPackageUpload.BodyReader(path)
	if err != nil {
		printer.Error(err)
	}

	printer.Info(fmt.Sprintf("Uploading '%s' to work package ...", filepath.Base(path)))

	body := &requests.RequestData{ContentType: contentType, Body: reader}
	status, response := requests.Do(link.Method, link.Href, nil, body)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	} else {
		printer.Done()
	}
}

func action(workPackage *dtos.WorkPackageDto, action string) {
	foundAction := findAction(action, workPackage.Embedded.CustomActions)
	if foundAction == nil {
		printer.Info(fmt.Sprintf(
			"No unique available action from input '%s' found for work package [#%d]. Please use one of the actions listed below.",
			action,
			workPackage.Id,
		))
		availableActions := common.Reduce(
			workPackage.Embedded.CustomActions,
			func(acc []*models.CustomAction, dto *dtos.CustomActionDto) []*models.CustomAction {
				return append(acc, dto.Convert())
			},
			[]*models.CustomAction{},
		)
		printer.CustomActions(availableActions)
		return
	}

	executeAction(workPackage, foundAction)
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

func executeAction(workPackage *dtos.WorkPackageDto, action *dtos.CustomActionDto) {
	printer.Info(fmt.Sprintf("Executing action '%s' on work package [#%d] ...", action.Name, workPackage.Id))

	requestBody := &dtos.CustomActionExecuteDto{
		LockVersion: workPackage.LockVersion,
		Links:       &dtos.ExecuteLinksDto{WorkPackage: &dtos.LinkDto{Href: workPackage.Links.Self.Href}},
	}

	b, err := json.Marshal(requestBody)
	if err != nil {
		printer.Error(err)
	}

	body := &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(b)}
	status, response := requests.Do(action.Links.Execute.Method, action.Links.Execute.Href, nil, body)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	} else {
		printer.Done()
	}
}
