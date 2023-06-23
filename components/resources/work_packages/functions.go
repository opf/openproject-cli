package work_packages

import (
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
	"github.com/opf/openproject-cli/models"
)

const path = "api/v3/work_packages"

func Lookup(id int64) *models.WorkPackage {
	return fetch(id).convert()
}

func All() []*models.WorkPackage {
	return fetchAll().convert()
}

func Update(id int64, action string) {
	workPackage := fetch(id)

	foundAction := findAction(action, workPackage.Embeddded.CustomActions)
	if foundAction == nil {
		printer.Info(fmt.Sprintf(
			"No unique available action from input '%s' found for work package [#%d]. Please use one of the actions listed below.",
			action,
			id,
		))
		availableActions := common.Reduce(
			workPackage.Embeddded.CustomActions,
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

	body, err := json.Marshal(reqBody)
	if err != nil {
		printer.Error(err)
	}

	status, response := requests.Do(action.Links.Execute.Method, action.Links.Execute.Href, nil, body)
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

func fetchAll() *WorkPackageCollectionDto {
	status, response := requests.Get(path, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	element := parser.Parse[WorkPackageCollectionDto](response)
	return &element
}
