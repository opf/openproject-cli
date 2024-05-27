package filters

import (
	"time"

	"github.com/opf/openproject-cli/components/requests"
)

type TimestampFilter struct {
	value string
}

func (f *TimestampFilter) Value() *string {
	return &f.value
}

func (f *TimestampFilter) Name() string {
	return "timestamp"
}

func (f *TimestampFilter) ShortHand() string {
	return ""
}

func (f *TimestampFilter) Usage() string {
	return `Returns the list of work packages at a specific timestamp. The timestamp should
be in the format of 'YYYY-MM-DDTHH:MM:SSZ' or date only 'YYYY-MM-DD', which
then assumes the time being 00:00:00Z.`
}

func (f *TimestampFilter) ValidateInput() error {
	_, err := time.Parse(time.DateOnly, f.value)
	if err == nil {
		return nil
	}

	_, err = time.Parse(time.RFC3339, f.value)
	return err
}

func (f *TimestampFilter) DefaultValue() string {
	return ""
}

func (f *TimestampFilter) Query() requests.Query {
	return requests.NewQuery(map[string]string{"timestamps": f.value}, nil)
}

func NewTimestampFilter() *TimestampFilter {
	return &TimestampFilter{}
}
