package notifications

import "github.com/opf/openproject-cli/components/resources"

type NotificationDto struct {
	Id        int64                `json:"id"`
	Reason    string               `json:"reason"`
	ReadIAN   bool                 `json:"readIAN"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
	Links     NotificationLinksDto `json:"_links"`
}

type NotificationLinksDto struct {
	Resource resources.LinkDto `json:"resource"`
}

type elements struct {
	Elements []*NotificationDto `json:"elements"`
}

type NotificationCollectionDto struct {
	Embedded elements `json:"_embedded"`
	Type     string   `json:"_type"`
}
