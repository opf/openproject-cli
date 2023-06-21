package requests_test

import (
	"testing"

	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/models/types"
)

func TestQuery_String_WithFilters(t *testing.T) {
	expectedString := "pageSize=100&filters=%5B%7B%22readIAN%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%22f%22%5D%7D%7D%2C%7B%22reason%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%22assigned%22%5D%7D%7D%5D"
	
	queryString := requests.NewQuery([]requests.Filter{
		requests.NewNotificationReadFilter(false),
		requests.NewNotificationReasonFilter(types.Assigned),
	}).String()

	if queryString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, queryString)
	}
}

func TestQuery_String_WithoutFilters(t *testing.T) {
	expectedString := "pageSize=100"
	
	queryString := requests.NewQuery([]requests.Filter{}).String()

	if queryString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, queryString)
	}
}
