package requests_test

import (
	"github.com/opf/openproject-cli/components/resources/notifications"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/requests"
)

func TestFilterQuery_String_WithFilters(t *testing.T) {
	expectedString := "filters=%5B%7B%22readIAN%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%22f%22%5D%7D%7D%2C%7B%22reason%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%22assigned%22%5D%7D%7D%5D&pageSize=100"

	queryString := requests.NewFilterQuery([]requests.Filter{
		notifications.ReadFilter(false),
		notifications.ReasonFilter("assigned"),
	}).String()

	if queryString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, queryString)
	}
}

func TestFilterQuery_String_WithoutFilters(t *testing.T) {
	expectedString := "pageSize=100"

	queryString := requests.NewFilterQuery([]requests.Filter{}).String()

	if queryString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, queryString)
	}
}

func TestPaginatedQuery_String_WithFilters(t *testing.T) {
	expectedString := "filters=%5B%7B%22status%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%221%22%2C%223%22%5D%7D%7D%5D&pageSize=7"

	filters := []requests.Filter{
		work_packages.StatusFilter("1,3"),
	}
	queryString := requests.NewPaginatedQuery(7, filters).String()

	if queryString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, queryString)
	}
}

func TestNewUnpaginatedQuery_String_WithFilters(t *testing.T) {
	expectedString := "filters=%5B%7B%22status%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%2212%22%2C%223%22%5D%7D%7D%5D&pageSize=-1"

	filters := []requests.Filter{
		work_packages.StatusFilter("12,3"),
	}
	queryString := requests.NewUnpaginatedQuery(nil, filters).String()

	if queryString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, queryString)
	}
}

func TestQuery_String(t *testing.T) {
	attributes := map[string]string{
		"pageSize":           "20",
		"includeSubprojects": "true",
		"timestamps":         "PT0S",
	}

	filters := []requests.Filter{
		work_packages.StatusFilter("1,3"),
		work_packages.TypeFilter("!1"),
	}
	queryString := requests.NewQuery(attributes, filters).String()

	if strings.Count(queryString, "&") != 3 {
		t.Errorf("Expected %s to contain 3 delimiter '&'.", queryString)
	}

	expected := "pageSize=20"
	contained := strings.Contains(queryString, expected)
	if !contained {
		t.Errorf("Expected %s to contain %s", queryString, expected)
	}

	expected = "includeSubprojects=true"
	contained = strings.Contains(queryString, expected)
	if !contained {
		t.Errorf("Expected %s to contain %s", queryString, expected)
	}

	expected = "timestamps=PT0S"
	contained = strings.Contains(queryString, expected)
	if !contained {
		t.Errorf("Expected %s to contain %s", queryString, expected)
	}

	expected = "filters=%5B%7B%22status%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%221%22%2C%223%22%5D%7D%7D%2C%7B%22type%22%3A%7B%22operator%22%3A%22%21%22%2C%22values%22%3A%5B%221%22%5D%7D%7D%5D"
	contained = strings.Contains(queryString, expected)
	if !contained {
		t.Errorf("Expected %s to contain %s", queryString, expected)
	}
}

func TestQuery_Merge(t *testing.T) {
	attributes1 := map[string]string{
		"pageSize":   "20",
		"timestamps": "PT0S",
	}
	filters1 := []requests.Filter{
		work_packages.StatusFilter("1,3"),
	}

	attributes2 := map[string]string{
		"pageSize":           "25",
		"includeSubprojects": "true",
	}
	filters2 := []requests.Filter{
		work_packages.TypeFilter("!1"),
	}

	query1 := requests.NewQuery(attributes1, filters1)
	query2 := requests.NewQuery(attributes2, filters2)

	result := query1.Merge(query2)
	expected := "filters=%5B%7B%22status%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%221%22%2C%223%22%5D%7D%7D%2C%7B%22type%22%3A%7B%22operator%22%3A%22%21%22%2C%22values%22%3A%5B%221%22%5D%7D%7D%5D&pageSize=25&timestamps=PT0S&includeSubprojects=true"

	if result.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, result.String())
	}
}
