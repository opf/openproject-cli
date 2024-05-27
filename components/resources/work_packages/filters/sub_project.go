package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

type SubProjectFilter struct {
	value string
}

func (f *SubProjectFilter) ValuePointer() *string {
	return &f.value
}

func (f *SubProjectFilter) Value() string {
	return f.value
}

func (f *SubProjectFilter) Name() string {
	return "sub-project"
}

func (f *SubProjectFilter) ShortHand() string {
	return ""
}

func (f *SubProjectFilter) Usage() string {
	return `Show only work packages of the specified subprojects. This filter only applies,
if the flag '--include-sub-projects' is set. It then includes only the sub
projects matching the filter. The value can be a single ID or a comma separated
array of IDs, i.e. '7,13'. Multiple values are concatenated with a logical 'OR'.`
}

func (f *SubProjectFilter) ValidateInput() error {
	matched, _ := regexp.Match("^([0-9,]+)$", []byte(f.value))
	if !matched {
		return errors.Custom(fmt.Sprintf("Invalid sub-project value %s.", printer.Yellow(f.value)))
	}

	return nil
}

func (f *SubProjectFilter) DefaultValue() string {
	return ""
}

func (f *SubProjectFilter) Query() requests.Query {
	return requests.NewQuery(nil, []requests.Filter{
		{
			Operator: "=",
			Name:     "subprojectId",
			Values:   strings.Split(f.value, ","),
		},
	})
}

func NewSubProjectFilter() *SubProjectFilter {
	return &SubProjectFilter{}
}
