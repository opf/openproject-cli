package dtos_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/dtos"
)

// A freshly created work package has lockVersion 0. The PATCH body must
// still contain it: omitting lockVersion makes the API respond with
// 409 UpdateConflict on every update of a fresh work package.
func TestWorkPackageDtoMarshalKeepsZeroLockVersion(t *testing.T) {
	patch := dtos.WorkPackageDto{LockVersion: 0}

	data, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	if !strings.Contains(string(data), `"lockVersion":0`) {
		t.Errorf("expected marshaled patch to contain \"lockVersion\":0, got %s", string(data))
	}
}
