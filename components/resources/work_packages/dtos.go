package work_packages

import (
	res "github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/components/resources/custom_actions"
)

type workPackageLinksDto struct {
	Self          *res.LinkDto    `json:"self"`
	Status        *res.LinkDto    `json:"status"`
	Assignee      *res.LinkDto    `json:"assignee"`
	Type          *res.LinkDto    `json:"type"`
	CustomActions []*res.LinkDto `json:"customActions"`
}

type workPackageDescription struct {
	Raw string `json:"raw"`
}

type WorkPackageDto struct {
	Id          int64                  `json:"id"`
	Subject     string                 `json:"subject"`
	Links       *workPackageLinksDto    `json:"_links"`
	Description *workPackageDescription `json:"description"`
	Embeddded   *embeddedDto            `json:"_embedded"`
	LockVersion int                    `json:"lockVersion"`
}

type embeddedDto struct {
	CustomActions []*custom_actions.CustomActionDto `json:"customActions"`
}

type elements struct {
	Elements []*WorkPackageDto `json:"elements"`
}

type WorkPackageCollectionDto struct {
	Embedded elements `json:"_embedded"`
	Type     string   `json:"_type"`
}
