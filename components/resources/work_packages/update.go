package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/opf/openproject-cli/components/common"
	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type UpdateOption int

const (
	UpdateCustomAction UpdateOption = iota
	UpdateAssignee
	UpdateAttachment
	UpdateDescription
	UpdateSubject
	UpdateType
)

var patchableUpdates = []UpdateOption{UpdateSubject, UpdateType, UpdateAssignee, UpdateDescription}

var patchMap = map[UpdateOption]func(patch, workPackage *dtos.WorkPackageDto, input string) (string, error){
	UpdateAssignee:    assigneePatch,
	UpdateDescription: descriptionPatch,
	UpdateType:        typePatch,
	UpdateSubject:     subjectPatch,
}

func Update(id string, options map[UpdateOption]string) (*models.WorkPackage, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	var customAction *dtos.CustomActionDto
	if input, ok := options[UpdateCustomAction]; ok {
		customAction, err = resolveAction(workPackage, input)
		if err != nil {
			return nil, err
		}
	}

	patchDto, updateString, patchNeeded, err := preparePatch(workPackage, options)
	if err != nil {
		return nil, err
	}

	if file, ok := options[UpdateAttachment]; ok {
		if err := validateAttachment(workPackage, file); err != nil {
			return nil, err
		}
	}

	if customAction != nil {
		err = executeAction(workPackage, customAction)
		if err != nil {
			return nil, err
		}

		// reload work package to get new lock version
		workPackage, err = fetch(id)
		if err != nil {
			return nil, err
		}
	}

	if patchNeeded {
		patchDto.LockVersion = workPackage.LockVersion
		if err := executePatch(workPackage, patchDto, updateString); err != nil {
			return nil, err
		}
	}

	if file, ok := options[UpdateAttachment]; ok {
		err = upload(workPackage, file)
		if err != nil {
			return nil, err
		}
	}

	workPackage, err = fetch(id)
	if err != nil {
		return nil, err
	}

	return workPackage.Convert(), nil
}

func preparePatch(workPackage *dtos.WorkPackageDto, options map[UpdateOption]string) (*dtos.WorkPackageDto, string, bool, error) {
	var patchNeeded = false
	patchDto := &dtos.WorkPackageDto{}
	var updateString string

	for option, value := range options {
		if !common.Contains(patchableUpdates, option) {
			continue
		}

		patchNeeded = true
		updateStringLine, err := patchMap[option](patchDto, workPackage, value)
		if err != nil {
			return nil, "", false, err
		}

		if len(updateStringLine) > 0 {
			if len(updateString) > 0 {
				updateString += "\n"
			}
			updateString += fmt.Sprintf("\t%s", updateStringLine)
		}
	}

	return patchDto, updateString, patchNeeded, nil
}

func executePatch(workPackage, patchDto *dtos.WorkPackageDto, updateString string) error {
	printer.Info("Updating work package with patch ...")

	marshal, err := json.Marshal(patchDto)
	if err != nil {
		return err
	}

	_, err = requests.Patch(workPackage.Links.Self.Href, &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(marshal)})
	if err != nil {
		return err
	}

	printer.Info(updateString)
	printer.Done()
	return nil
}

func validateAttachment(workPackage *dtos.WorkPackageDto, path string) error {
	if workPackage.Links == nil {
		return fmt.Errorf("this work package does not accept attachments (missing permission?)")
	}
	if workPackage.Links.PrepareAttachment != nil {
		return fmt.Errorf("uploads to fog storages are currently not supported")
	}
	if workPackage.Links.AddAttachment == nil {
		return fmt.Errorf("this work package does not accept attachments (missing permission?)")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("attachment path %q is a directory", path)
	}
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

		return "", openerrors.ErrHandled
	}

	if patch.Links == nil {
		patch.Links = &dtos.WorkPackageLinksDto{}
	}

	patch.Links.Type = foundType.Links.Self
	return fmt.Sprintf("Type -> %s", foundType.Name), nil
}

func subjectPatch(patch, _ *dtos.WorkPackageDto, input string) (string, error) {
	patch.Subject = input
	return fmt.Sprintf("Subject -> %s", input), nil
}

func descriptionPatch(patch, _ *dtos.WorkPackageDto, input string) (string, error) {
	patch.Description = &dtos.LongTextDto{Format: "markdown", Raw: input}
	return "Description updated", nil
}

func assigneePatch(patch, _ *dtos.WorkPackageDto, input string) (string, error) {
	userId, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid user id %q: must be a number", input)
	}

	if patch.Links == nil {
		patch.Links = &dtos.WorkPackageLinksDto{}
	}

	patch.Links.Assignee = &dtos.LinkDto{Href: paths.User(userId)}
	return fmt.Sprintf("Assignee -> %s", input), nil
}
