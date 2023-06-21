package notifications

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/models"
)

const path = "api/v3/notifications"

func All() []*models.Notification {
	status, response := requests.Get(path)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}
	
	element := parser.Parse[NotificationCollectionDto](response)
	return element.convert()
}
