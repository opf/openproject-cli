package routes

import (
	"fmt"
	"net/url"

	"github.com/opf/openproject-cli/models"
)

var host *url.URL

func Init(h *url.URL) {
	host = h
}

func WorkPackageUrl(workPackage *models.WorkPackage) url.URL {
	routeUrl := *host
	routeUrl.Path = fmt.Sprintf("work_packages/%d", workPackage.Id)
	return routeUrl
}

func ProjectUrl(project *models.Project) url.URL {
	routeUrl := *host
	routeUrl.Path = fmt.Sprintf("projects/%d", project.Id)
	return routeUrl
}
