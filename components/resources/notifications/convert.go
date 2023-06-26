package notifications

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/models"
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
		Id:              uint64(dto.Id),
		ResourceId:      parser.IdFromLink(dto.Links.Resource.Href),
		ResourceSubject: dto.Links.Resource.Title,
		Reason:          dto.Reason,
		Read:            dto.ReadIAN,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
