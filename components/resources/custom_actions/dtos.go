package custom_actions

import "github.com/opf/openproject-cli/components/resources"

type CustomActionDto struct {
	Name  string    `json:"name"`
	Links *linksDto `json:"_links"`
}

type linksDto struct {
	Self    *resources.LinkDto `json:"self"`
	Execute *resources.LinkDto `json:"executeImmediately"`
}

type CustomActionExecuteDto struct {
	LockVersion int              `json:"lockVersion"`
	Links       *ExecuteLinksDto `json:"_links"`
}

type ExecuteLinksDto struct {
	WorkPackage *resources.LinkDto `json:"workPackage"`
}
