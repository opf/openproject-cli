package projects_test

import (
	"testing"

	"github.com/opf/openproject-cli/components/resources/projects"
)

func TestValidateIdentifier(t *testing.T) {
	valid := []string{
		"devops",
		"my-project",
		"project_name",
		"42",
		"ABC123",
	}

	for _, id := range valid {
		if err := projects.ValidateIdentifier(id); err != nil {
			t.Errorf("ValidateIdentifier(%q) returned unexpected error: %v", id, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"", "empty string"},
		{"my project", "space"},
		{"project/path", "slash"},
		{"proj@name", "at sign"},
		{"proj.name", "dot"},
		{"proj!name", "exclamation mark"},
		{"project+extra", "plus sign"},
	}

	for _, tc := range invalid {
		if err := projects.ValidateIdentifier(tc.input); err == nil {
			t.Errorf("ValidateIdentifier(%q) (%s) expected error but got nil", tc.input, tc.desc)
		}
	}
}
