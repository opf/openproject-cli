package work_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
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

var patchableUpdates = []UpdateOption{Subject, Type}

var patchMap = map[UpdateOption]func(patch, workPackage *dtos.WorkPackageDto, input string) (string, error){
	Type:    typePatch,
	Subject: subjectPatch,
}

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
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	if customAction, ok := options[Action]; ok {
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

	if file, ok := options[Attach]; ok {
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
	types, err := availableTypes(workPackage)
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

		printer.Types(common.Reduce(types,
			func(acc []*models.Type, dto *dtos.TypeDto) []*models.Type {
				return append(acc, dto.Convert())
			}, []*models.Type{}))
		return "", nil
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

func fetch(id uint64) (*dtos.WorkPackageDto, error) {
	response, err := requests.Get(paths.WorkPackage(id), nil)
	if err != nil {
		return nil, err
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)
	return &workPackage, nil
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
