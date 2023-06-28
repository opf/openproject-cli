package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestNotification(t *testing.T) {
	testingPrinter.Reset()

	notification := &models.Notification{
		ResourceId:      15,
		ResourceSubject: "Subject 1",
		Reason:          "Comment created",
	}

	idString := "#" + strconv.FormatUint(notification.ResourceId, 10)
	expected := fmt.Sprintf("%s (%s) %s\n", printer.Red(idString), notification.Reason, printer.Cyan(notification.ResourceSubject))

	printer.Notification(notification)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}

func TestNotifications(t *testing.T) {
	testingPrinter.Reset()

	notifications := []*models.Notification{
		{
			Id:              1,
			ResourceId:      151,
			ResourceSubject: "Subject 1",
			Reason:          "Comment created",
		},
		{
			Id:              2,
			ResourceId:      151,
			ResourceSubject: "Subject 1",
			Reason:          "Comment created",
		},
		{
			Id:              3,
			ResourceId:      12,
			ResourceSubject: "Subject 2",
			Reason:          "Status changed",
		},
	}

	var expected string
	expected += fmt.Sprintf("%s %s %s %s\n", printer.Red(fmt.Sprintf("#%d", notifications[0].ResourceId)), fmt.Sprintf("(%s)", notifications[0].Reason), printer.Cyan(notifications[0].ResourceSubject), printer.Green(fmt.Sprintf("(2)")))
	expected += fmt.Sprintf("%s %s %s\n", printer.Red(fmt.Sprintf(" #%d", notifications[2].ResourceId)), fmt.Sprintf("(%s) ", notifications[2].Reason), printer.Cyan(notifications[2].ResourceSubject))

	printer.Notifications(notifications)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
