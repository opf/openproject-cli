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

func TestPagedQuery_String_WithFilters(t *testing.T) {
	expectedString := "filters=%5B%7B%22status%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%221%22%2C%223%22%5D%7D%7D%5D&pageSize=-1"

	filters := []requests.Filter{
		work_packages.StatusFilter("1,3"),
	}
	queryString := requests.NewPagedQuery(-1, filters).String()

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
