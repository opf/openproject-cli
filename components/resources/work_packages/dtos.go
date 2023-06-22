package work_packages

import "github.com/opf/openproject-cli/components/resources"

type workPackageLinksDto struct {
	Status   resources.LinkDto `json:"status"`
	Assignee resources.LinkDto `json:"assignee"`
	Type     resources.LinkDto `json:"type"`
}

type workPackageDescription struct {
	Raw string `json:"raw"`
}

type WorkPackageDto struct {
	Id          int64                  `json:"id"`
	Subject     string                 `json:"subject"`
	Links       workPackageLinksDto    `json:"_links"`
	Description workPackageDescription `json:"description"`
}
