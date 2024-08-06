package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

type VersionFilter struct {
	value string
}

func (f *VersionFilter) ValuePointer() *string {
	return &f.value
}

func (f *VersionFilter) Value() string {
	return f.value
}

func (f *VersionFilter) Name() string {
	return "version"
}

func (f *VersionFilter) ShortHand() string {
	return "v"
}

func (f *VersionFilter) Usage() string {
	return `Show only work packages that are assigned to the given version. The value can be
a single ID or a comma separated array of IDs, i.e. '7,13'.`
}

func (f *VersionFilter) ValidateInput() error {
	matched, _ := regexp.Match("^([0-9,]+)$", []byte(f.value))
	if !matched {
		return errors.Custom(fmt.Sprintf("Invalid version value %s.", printer.Yellow(f.value)))
	}

	return nil
}

func (f *VersionFilter) DefaultValue() string {
	return ""
}

func (f *VersionFilter) Query() requests.Query {
	return requests.NewQuery(nil, []requests.Filter{
		{
			Operator: "=",
			Name:     "version",
			Values:   strings.Split(f.value, ","),
		},
	})
}

func NewVersionFilter() *VersionFilter {
	return &VersionFilter{}
}
