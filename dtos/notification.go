package dtos

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/models"
)

type NotificationDto struct {
	Id        int64                `json:"id"`
	Reason    string               `json:"reason"`
	ReadIAN   bool                 `json:"readIAN"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
	Links     NotificationLinksDto `json:"_links"`
}

type NotificationLinksDto struct {
	Resource LinkDto `json:"resource"`
}

type notificationElements struct {
	Elements []*NotificationDto `json:"elements"`
}

type NotificationCollectionDto struct {
	Embedded *notificationElements `json:"_embedded"`
	Type     string                `json:"_type"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *NotificationCollectionDto) Convert() []*models.Notification {
	var notifications = make([]*models.Notification, len(dto.Embedded.Elements))

	for idx, n := range dto.Embedded.Elements {
		notifications[idx] = n.Convert()
	}

	return notifications
}

func (dto *NotificationDto) Convert() *models.Notification {
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