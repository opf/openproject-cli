package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestNotifications(t *testing.T) {
	testingPrinter.Reset()

	notifications := []*models.Notification{
		{
			Id:              42,
			ResourceId:      15,
			ResourceSubject: "Subject 1",
			Reason:          "Comment created",
			Read:            false,
			CreatedAt:       "",
			UpdatedAt:       "",
		},
		{
			Id:              54,
			ResourceId:      12,
			ResourceSubject: "Subject 2",
			Reason:          "Status changed",
			Read:            false,
			CreatedAt:       "",
			UpdatedAt:       "",
		},
	}

	expected := common.Reduce[*models.Notification, string](
		notifications,
		func(state string, notification *models.Notification) string {
			idString := "#" + strconv.FormatUint(notification.ResourceId, 10)
			return state + fmt.Sprintf("[%s] %s (%s)\n", printer.Red(idString), printer.Cyan(notification.ResourceSubject), notification.Reason)
		},
		"")

	printer.Notifications(notifications)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
