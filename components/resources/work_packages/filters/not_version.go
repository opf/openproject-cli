package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

type NotVersionFilter struct {
	value string
}

func (f *NotVersionFilter) ValuePointer() *string {
	return &f.value
}

func (f *NotVersionFilter) Value() string {
	return f.value
}

func (f *NotVersionFilter) Name() string {
	return "not-version"
}

func (f *NotVersionFilter) ShortHand() string {
	return ""
}

func (f *NotVersionFilter) Usage() string {
	return `Show only work packages that are not assigned to the given version. The value
can be a single ID or a comma separated array of IDs, i.e. '7,13'.`
}

func (f *NotVersionFilter) ValidateInput() error {
	matched, _ := regexp.Match("^([0-9,]+)$", []byte(f.value))
	if !matched {
		return errors.Custom(fmt.Sprintf("Invalid version value %s.", printer.Yellow(f.value)))
	}

	return nil
}

func (f *NotVersionFilter) DefaultValue() string {
	return ""
}

func (f *NotVersionFilter) Query() requests.Query {
	return requests.NewQuery(nil, []requests.Filter{
		{
			Operator: "!",
			Name:     "version",
			Values:   strings.Split(f.value, ","),
		},
	})
}

func NewNotVersionFilter() *NotVersionFilter {
	return &NotVersionFilter{}
}
