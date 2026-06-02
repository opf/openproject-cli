package work_packages_test

import (
	"testing"

	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestValidateIdentifier(t *testing.T) {
	valid := []string{
		"1",
		"72427",
		"SJF-13",
		"A-1",
		"MYPROJECT-100",
		"MY_PROJ-5",
		"ABCDEFGHIJ-1",
	}

	for _, id := range valid {
		if err := work_packages.ValidateIdentifier(id); err != nil {
			t.Errorf("ValidateIdentifier(%q) returned unexpected error: %v", id, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"", "empty string"},
		{"abc", "lowercase letters only"},
		{"sjf-13", "lowercase project identifier"},
		{"1234-5", "project identifier starts with digit"},
		{"SJF-abc", "non-numeric sequence number"},
		{"SJF-", "missing sequence number"},
		{"-13", "missing project identifier"},
		{"SJF_13", "underscore instead of hyphen separator"},
		{"ABCDEFGHIJK-1", "project identifier exceeds 10 characters"},
		{"SJF 13", "space"},
	}

	for _, tc := range invalid {
		if err := work_packages.ValidateIdentifier(tc.input); err == nil {
			t.Errorf("ValidateIdentifier(%q) (%s) expected error but got nil", tc.input, tc.desc)
		}
	}
}
