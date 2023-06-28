package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

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
