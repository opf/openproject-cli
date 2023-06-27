package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func Notifications(v interface{}) {
	list, ok := v.([]*models.Notification)
	if ok {
		for _, n := range list {
			printNotification(n)
		}
	}

	single, ok := v.(*models.Notification)
	if ok {
		printNotification(single)
	}
}

func printNotification(n *models.Notification) {
	id := fmt.Sprintf("#%d", n.ResourceId)
	activePrinter.Printf("[%s] %s (%s)\n", Red(id), Cyan(n.ResourceSubject), n.Reason)
}
