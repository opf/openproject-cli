package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/status"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

type UpdateOption int

const (
	UpdateCustomAction UpdateOption = iota
	UpdateAssignee
	UpdateAttachment
	UpdateSubject
	UpdateDescription
	UpdateType
	UpdateStatus
)

var patchableUpdates = []UpdateOption{UpdateSubject, UpdateType, UpdateAssignee, UpdateDescription, UpdateStatus}

var patchMap = map[UpdateOption]func(patch, workPackage *dtos.WorkPackageDto, input string) error{
	UpdateAssignee:    assigneePatch,
	UpdateType:        typePatch,
	UpdateSubject:     subjectPatch,
	UpdateDescription: descriptionPatch,
	UpdateStatus:      statusPatch,
}

func DryRunUpdate(id uint64, options map[UpdateOption]string) (*models.WorkPackageUpdatePlan, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	plan := &models.WorkPackageUpdatePlan{
		Valid:          true,
		Operation:      "update",
		WorkPackageID:  id,
		Subject:        options[UpdateSubject],
		Status:         options[UpdateStatus],
		Action:         options[UpdateCustomAction],
		Attach:         options[UpdateAttachment],
		ResolvedFields: map[string]models.ResolvedField{},
	}

	if description, ok := options[UpdateDescription]; ok {
		plan.Description = &description
	}

	if assignee, ok := options[UpdateAssignee]; ok {
		plan.Assignee = assignee
	}

	if value, ok := options[UpdateStatus]; ok {
		resolvedStatus, err := resolveStatus(value)
		if err != nil {
			return nil, err
		}

		plan.Status = resolvedStatus.Name
	}

	if value, ok := options[UpdateType]; ok {
		types, err := availableTypes(workPackage.Links.Project)
		if err != nil {
			return nil, err
		}

		foundType := findType(value, types)
		if foundType == nil {
			return nil, fmt.Errorf("no unique available type from input %q found for work package #%d", value, id)
		}

		plan.Type = foundType.Name
	}

	return plan, nil
}

func Update(id uint64, options map[UpdateOption]string) (*models.WorkPackage, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	if customAction, ok := options[UpdateCustomAction]; ok {
		err = action(workPackage, customAction)
		if err != nil {
			return nil, err
		}

		// reload work package to get new lock version
		workPackage, err = fetch(id)
		if err != nil {
			return nil, err
		}
	}

	err = patch(workPackage, options)
	if err != nil {
		return nil, err
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

func patch(workPackage *dtos.WorkPackageDto, options map[UpdateOption]string) error {
	var patchNeeded = false
	patchDto := dtos.WorkPackageDto{LockVersion: workPackage.LockVersion}

	for _, option := range patchableUpdates {
		value, ok := options[option]
		if !ok {
			continue
		}

		patchNeeded = true
		if err := patchMap[option](&patchDto, workPackage, value); err != nil {
			return err
		}
	}

	if !patchNeeded {
		return nil
	}

	marshal, err := json.Marshal(patchDto)
	if err != nil {
		return err
	}

	_, err = requests.Patch(workPackage.Links.Self.Href, &requests.RequestData{ContentType: "application/json", Body: bytes.NewReader(marshal)})
	if err != nil {
		return err
	}
	return nil
}

func typePatch(patch, workPackage *dtos.WorkPackageDto, input string) error {
	types, err := availableTypes(workPackage.Links.Project)
	if err != nil {
		return err
	}

	foundType := findType(input, types)
	if foundType == nil {
		return fmt.Errorf("no unique available type from input %q found for project #%d", input, parser.IdFromLink(workPackage.Links.Project.Href))
	}

	if patch.Links == nil {
		patch.Links = &dtos.WorkPackageLinksDto{}
	}

	patch.Links.Type = foundType.Links.Self
	return nil
}

func subjectPatch(patch, _ *dtos.WorkPackageDto, input string) error {
	patch.Subject = input
	return nil
}

func assigneePatch(patch, _ *dtos.WorkPackageDto, input string) error {
	userId, _ := strconv.ParseUint(input, 10, 64)

	if patch.Links == nil {
		patch.Links = &dtos.WorkPackageLinksDto{}
	}

	patch.Links.Assignee = &dtos.LinkDto{Href: paths.User(userId)}
	return nil
}

func descriptionPatch(patch, _ *dtos.WorkPackageDto, input string) error {
	patch.Description = &dtos.LongTextDto{Raw: input}
	return nil
}

func statusPatch(patch, _ *dtos.WorkPackageDto, input string) error {
	resolvedStatus, err := resolveStatus(input)
	if err != nil {
		return err
	}

	if patch.Links == nil {
		patch.Links = &dtos.WorkPackageLinksDto{}
	}

	patch.Links.Status = &dtos.LinkDto{Href: paths.StatusById(resolvedStatus.Id)}
	return nil
}

func resolveStatus(input string) (*models.Status, error) {
	statuses, err := status.All()
	if err != nil {
		return nil, err
	}

	for _, candidate := range statuses {
		if strings.EqualFold(candidate.Name, input) {
			return candidate, nil
		}
	}

	return nil, fmt.Errorf("no status named %q found", input)
}
