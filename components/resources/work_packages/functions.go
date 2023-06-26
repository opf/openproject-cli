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
	"github.com/opf/openproject-cli/components/resources"
	actions "github.com/opf/openproject-cli/components/resources/custom_actions"
	workPackageUpload "github.com/opf/openproject-cli/components/resources/work_packages/upload"
	"github.com/opf/openproject-cli/models"
)

type UpdateOption int

const (
	Action UpdateOption = iota
	Attach
)

const path = "api/v3/work_packages"

func Lookup(id int64) *models.WorkPackage {
	return fetch(id).convert()
}

func All(principal *models.Principal) []*models.WorkPackage {
	var filters []requests.Filter

	if principal != nil {
		filters = append(filters, AssigneeFilter(principal))
	}

	query := requests.NewQuery(filters)

	status, response := requests.Get(path, &query)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	element := parser.Parse[WorkPackageCollectionDto](response)
	return element.convert()
}

func Update(id int64, opts map[UpdateOption]string) {
	workPackage := fetch(id)

	for updateOpt, value := range opts {
		switch updateOpt {
		case Action:
			action(workPackage, value)
		case Attach:
			upload(workPackage, value)
		}
	}
}

func upload(dto *WorkPackageDto, path string) {
	if dto.Links.PrepareAttachment != nil {
		printer.ErrorText(fmt.Sprintf("Uploads to fog storages are currently not supported. :("))
	}

	link := dto.Links.AddAttachment
	reader, contentType, err := workPackageUpload.BodyReader(path)
	if err != nil {
		printer.Error(err)
	}

	printer.Info(fmt.Sprintf("Uploading '%s' to work package ...", filepath.Base(path)))

	body := &requests.RequestBody{ContentType: contentType, Body: reader}
	status, response := requests.Do(link.Method, link.Href, nil, body)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	} else {
		printer.Done()
	}
}

func fetch(id int64) *WorkPackageDto {
	status, response := requests.Get(filepath.Join(path, strconv.FormatInt(id, 10)), nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	workPackage := parser.Parse[WorkPackageDto](response)
	return &workPackage
}

func action(workPackage *WorkPackageDto, action string) {
	foundAction := findAction(action, workPackage.Embedded.CustomActions)
	if foundAction == nil {
		printer.Info(fmt.Sprintf(
			"No unique available action from input '%s' found for work package [#%d]. Please use one of the actions listed below.",
			action,
			workPackage.Id,
		))
		availableActions := common.Reduce(
			workPackage.Embedded.CustomActions,
			func(acc []*models.CustomAction, dto *actions.CustomActionDto) []*models.CustomAction {
				return append(acc, dto.Convert())
			},
			[]*models.CustomAction{},
		)
		printer.CustomActions(availableActions)
		return
	}

	executeAction(workPackage, foundAction)
}

func findAction(actionInput string, availableActions []*actions.CustomActionDto) *actions.CustomActionDto {
	var actionAsId = false
	actionId, err := strconv.ParseInt(actionInput, 10, 64)
	if err == nil {
		actionAsId = true
	}

	var found []*actions.CustomActionDto
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

func executeAction(workPackage *WorkPackageDto, action *actions.CustomActionDto) {
	printer.Info(fmt.Sprintf("Executing action '%s' on work package [#%d] ...", action.Name, workPackage.Id))

	reqBody := &actions.CustomActionExecuteDto{
		LockVersion: workPackage.LockVersion,
		Links:       &actions.ExecuteLinksDto{WorkPackage: &resources.LinkDto{Href: workPackage.Links.Self.Href}},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		printer.Error(err)
	}

	body := &requests.RequestBody{ContentType: "application/json", Body: bytes.NewReader(b)}
	status, response := requests.Do(action.Links.Execute.Method, action.Links.Execute.Href, nil, body)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	} else {
		printer.Done()
	}
}
