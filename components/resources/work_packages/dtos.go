package work_packages

type linkDto struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type linksDto struct {
	Status   linkDto `json:"status"`
	Assignee linkDto `json:"assignee"`
	Type     linkDto `json:"type"`
}

type WorkPackageDto struct {
	Id      int64    `json:"id"`
	Subject string   `json:"subject"`
	Links   linksDto `json:"_links"`
}
