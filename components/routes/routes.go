package routes

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/opf/openproject-cli/models"
)

var host *url.URL

func Init(h *url.URL) {
	host = h
}

func WorkPackageUrl(workPackage *models.WorkPackage) *url.URL {
	routeUrl := *host
	// Current servers always set displayId (`identifier.presence || id`), but
	// older versions omit it, so fall back to the numeric id.
	displayId := workPackage.DisplayId
	if displayId == "" {
		displayId = strconv.FormatUint(workPackage.Id, 10)
	}
	routeUrl.Path = fmt.Sprintf("wp/%s", displayId)
	return &routeUrl
}

func ProjectUrl(project *models.Project) *url.URL {
	routeUrl := *host
	routeUrl.Path = fmt.Sprintf("projects/%s", project.Identifier)
	return &routeUrl
}
