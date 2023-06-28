package notifications

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

const path = "api/v3/notifications"

func All(reason string) ([]*models.Notification, error) {
	response, err := requests.Get(path, generateQuery(reason))
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.NotificationCollectionDto](response)
	return element.Convert(), nil
}

func generateQuery(reason string) *requests.Query {
	filters := []requests.Filter{
		requests.NewNotificationReadFilter(false),
	}

	if reason != "" {
		filters = append(filters, requests.NewNotificationReasonFilter(reason))
	}

	query := requests.NewQuery(filters)
	return &query
}
