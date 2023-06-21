package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func Notifications(notifications []*models.Notification) {
	for _, notification := range notifications {
		printNotification(notification)
	}
}

func printNotification(n *models.Notification) {
	id := fmt.Sprintf("#%d", n.ResourceId)
	fmt.Printf("[%s] %s (%s)\n", red(id), cyan(n.ResourceSubject), n.Reason)
}
