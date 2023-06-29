package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type UpdateOption int

const (
	UpdateAction UpdateOption = iota
	UpdateAttach
	UpdateSubject
	UpdateType
)

var patchableUpdates = []UpdateOption{UpdateSubject, UpdateType}

var patchMap = map[UpdateOption]func(patch, workPackage *dtos.WorkPackageDto, input string) (string, error){
	UpdateType:    typePatch,
	UpdateSubject: subjectPatch,
}

func Update(id uint64, options map[UpdateOption]string) (*models.WorkPackage, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	if customAction, ok := options[UpdateAction]; ok {
		err = action(workPackage, customAction)
		if err != nil {
			printer.Error(err)
		} else {
			// reload work package to get new lock version
			workPackage, err = fetch(id)
			if err != nil {
				return nil, err
			}
		}
	}

	err = patch(workPackage, options)
	if err != nil {
		printer.Error(err)
	}

	if file, ok := options[UpdateAttach]; ok {
		err = upload(workPackage, file)
		if err != nil {
			printer.Error(err)
		}
	}

	workPackage, err = fetch(id)
	if err != nil {
		return nil, err
	}

	return workPackage.Convert(), nil
}

func patch(workPackage *dtos.WorkPackageDto, options map[UpdateOption]string) error {
	var patchNeeded = false
	patchDto := dtos.WorkPackageDto{LockVersion: workPackage.LockVersion}
	var updateString string

	for option, value := range options {
		if !common.Contains(patchableUpdates, option) {
			continue
		}

		patchNeeded = true
		updateStringLine, err := patchMap[option](&patchDto, workPackage, value)
		if err != nil {
			return err
		}

		if len(updateStringLine) > 0 {
			if len(updateString) > 0 {
				updateString += "\n"
			}
			updateString += fmt.Sprintf("\t%s", updateStringLine)
		}
	}

	if !patchNeeded {
		return nil
	}

	printer.Info(fmt.Sprintf("Updating work package with patch ..."))
	printer.Info(updateString)

	marshal, err := json.Marshal(patchDto)
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

func typePatch(patch, workPackage *dtos.WorkPackageDto, input string) (string, error) {
	types, err := availableTypes(workPackage.Links.Project)
	if err != nil {
		return "", err
	}

	foundType := findType(input, types)
	if foundType == nil {
		printer.ErrorText("Failed to update work package type.")
		printer.Info(fmt.Sprintf(
			"No unique available type from input %s found for project %s. Please use one of the types listed below.",
			printer.Cyan(input),
			printer.Red(fmt.Sprintf("#%d", parser.IdFromLink(workPackage.Links.Project.Href))),
		))

		printer.Types(types.Convert())

		return "", nil
	}

	if patch.Links == nil {
		patch.Links = &dtos.WorkPackageLinksDto{}
	}

	patch.Links.Type = foundType.Links.Self
	return fmt.Sprintf("UpdateType -> %s", foundType.Name), nil
}

func subjectPatch(patch, _ *dtos.WorkPackageDto, input string) (string, error) {
	patch.Subject = input
	return fmt.Sprintf("UpdateSubject -> %s", input), nil
}
