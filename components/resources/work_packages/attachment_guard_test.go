package work_packages

import (
	"testing"

	"github.com/opf/openproject-cli/dtos"
)

// A work package whose payload has no _links object must not panic the
// attachment validation; it should be reported as not accepting attachments.
func TestValidateAttachmentNilLinks(t *testing.T) {
	err := validateAttachment(&dtos.WorkPackageDto{Links: nil}, "irrelevant")
	if err == nil {
		t.Fatal("validateAttachment(nil Links) = nil, want error")
	}
}

// upload must self-guard: the work package may be re-fetched after a custom
// action, dropping the addAttachment link that validateAttachment saw earlier.
func TestUploadNilAddAttachment(t *testing.T) {
	cases := map[string]*dtos.WorkPackageDto{
		"nil links":          {Links: nil},
		"nil add attachment": {Links: &dtos.WorkPackageLinksDto{}},
	}
	for name, dto := range cases {
		t.Run(name, func(t *testing.T) {
			if err := upload(dto, "irrelevant"); err == nil {
				t.Errorf("upload(%s) = nil, want error", name)
			}
		})
	}
}
