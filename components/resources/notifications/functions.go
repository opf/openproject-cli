package notifications

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func All(reason string) ([]*models.Notification, error) {
	response, err := requests.Get(paths.Notifications(), generateQuery(reason))
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
