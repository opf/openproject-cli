package notifications

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/models"
	"github.com/opf/openproject-cli/models/types"
)

func (dto *NotificationCollectionDto) convert() []*models.Notification {
	var notifications = make([]*models.Notification, len(dto.Embedded.Elements))

	for idx, n := range dto.Embedded.Elements {
		notifications[idx] = n.convert()
	}

	return notifications
}

func (dto *NotificationDto) convert() *models.Notification {
	return &models.Notification{
		Id:              dto.Id,
		ResourceId:      parser.IdFromLink(dto.Links.Resource.Href),
		ResourceSubject: dto.Links.Resource.Title,
		Reason:          types.Reason(dto.Reason),
		Read:            dto.ReadIAN,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
