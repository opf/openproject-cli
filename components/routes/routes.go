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

func WorkPackage(workPackage *models.WorkPackage) string {
	routeUrl := *host
	routeUrl.Path = fmt.Sprintf("work_packages/%d", workPackage.Id)
	return routeUrl.String()
}
