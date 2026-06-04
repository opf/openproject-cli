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

func WorkPackageUrl(workPackage *models.WorkPackage) *url.URL {
	routeUrl := *host
	// DisplayId is always non-empty: the server code sets it to `identifier.presence || id`
	// (in WorkPackage::SemanticIdentifier), so it falls back to the numeric id automatically.
	routeUrl.Path = fmt.Sprintf("wp/%s", workPackage.DisplayId)
	return &routeUrl
}

func ProjectUrl(project *models.Project) *url.URL {
	routeUrl := *host
	routeUrl.Path = fmt.Sprintf("projects/%s", project.Identifier)
	return &routeUrl
}
