package notifications

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/models"
	"github.com/opf/openproject-cli/models/types"
)

const path = "api/v3/notifications"

func All(reason types.Reason) []*models.Notification {
	status, response := requests.Get(path, generateQuery(reason))
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	element := parser.Parse[NotificationCollectionDto](response)
	return element.convert()
}

func generateQuery(reason types.Reason) *requests.Query {
	filters := []requests.Filter{
		requests.NewNotificationReadFilter(false),
	}

	if reason != types.None {
		filters = append(filters, requests.NewNotificationReasonFilter(reason))
	}

	query := requests.NewQuery(filters)
	return &query
}
