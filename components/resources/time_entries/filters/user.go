package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

// validates the value is one or more of the following - separated by commas
// - a valid login (see https://github.com/opf/openproject/blob/dev/app/models/user.rb#L135)
// - a numeric id
// - me
const validValueRegexp = `^([\pL0-9_\-@.+ ]+|[0-9]+|me)(,([\pL0-9_\-@.+ ]+|[0-9]+|me))*$`

type UserFilter struct {
	value string
}

func (f *UserFilter) ValuePointer() *string {
	return &f.value
}

func (f *UserFilter) Value() string {
	return f.value
}

func (f *UserFilter) Name() string {
	return "user"
}

func (f *UserFilter) ShortHand() string {
	return "u"
}

func (f *UserFilter) Usage() string {
	return `User the time entry tracks expenditures for (can be name, ID or 'me')`
}

func (f *UserFilter) ValidateInput() error {
	matched, _ := regexp.Match(validValueRegexp, []byte(f.value))
	if !matched {
		return errors.Custom(fmt.Sprintf("Invalid user value %s.", printer.Yellow(f.value)))
	}

	return nil
}

func (f *UserFilter) DefaultValue() string {
	return "me"
}

func (f *UserFilter) Query() requests.Query {
	return requests.NewQuery(nil, []requests.Filter{
		{
			Operator: "=",
			Name:     "user",
			Values:   strings.Split(f.value, ","),
		},
	})
}

func NewUserFilter() *UserFilter {
	return &UserFilter{}
}
