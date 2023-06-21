package notifications

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

type LinkDto struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type elements struct {
	Elements []*NotificationDto `json:"elements"`
}

type NotificationCollectionDto struct {
	Embedded elements `json:"_embedded"`
	Type     string   `json:"_type"`
}
