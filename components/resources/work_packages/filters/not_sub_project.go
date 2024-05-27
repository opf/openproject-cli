package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

type NotSubProjectFilter struct {
	value string
}

func (f *NotSubProjectFilter) ValuePointer() *string {
	return &f.value
}

func (f *NotSubProjectFilter) Value() string {
	return f.value
}

func (f *NotSubProjectFilter) Name() string {
	return "not-sub-project"
}

func (f *NotSubProjectFilter) ShortHand() string {
	return ""
}

func (f *NotSubProjectFilter) Usage() string {
	return `Show only work packages that are not inside the specified subprojects. This
filter only applies, if the flag '--include-sub-projects' is set. It then
excludes all sub projects matching the filter. The value can be a single ID or
a comma separated array of IDs, i.e. '7,13'.`
}

func (f *NotSubProjectFilter) ValidateInput() error {
	matched, _ := regexp.Match("^([0-9,]+)$", []byte(f.value))
	if !matched {
		return errors.Custom(fmt.Sprintf("Invalid not-sub-project value %s.", printer.Yellow(f.value)))
	}

	return nil
}

func (f *NotSubProjectFilter) DefaultValue() string {
	return ""
}

func (f *NotSubProjectFilter) Query() requests.Query {
	return requests.NewQuery(nil, []requests.Filter{
		{
			Operator: "!",
			Name:     "subprojectId",
			Values:   strings.Split(f.value, ","),
		},
	})
}

func NewNotSubProjectFilter() *NotSubProjectFilter {
	return &NotSubProjectFilter{}
}
