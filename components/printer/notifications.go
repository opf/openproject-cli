package printer

import (
	"fmt"
	"strings"

	"github.com/opf/openproject-cli/models"
)

func Notification(n *models.Notification) {
	activeRenderer.Notification(n)
}

func Notifications(notifications []*models.Notification) {
	activeRenderer.Notifications(notifications)
}

type groupedNotification struct {
	notification *models.Notification
	count        int
}

func printNotification(line *groupedNotification, maxIdLength, maxReasonLength int) {
	var parts []string

	diff := maxIdLength - idLength(line.notification.ResourceId)
	idStr := fmt.Sprintf("%s#%d", indent(diff), line.notification.ResourceId)
	parts = append(parts, Red(idStr))

	diff = maxReasonLength - len(line.notification.Reason)
	typeStr := fmt.Sprintf("(%s)%s", line.notification.Reason, indent(diff))
	parts = append(parts, typeStr)

	parts = append(parts, Cyan(line.notification.ResourceSubject))

	if line.count > 1 {
		parts = append(parts, Green(fmt.Sprintf("(%d)", line.count)))
	}

	activePrinter.Println(strings.Join(parts, " "))
}

func group(notifications []*models.Notification) []*groupedNotification {
	var list []*groupedNotification
	for _, notification := range notifications {
		var alreadyAdded = false
		for _, grouped := range list {
			if isSameNotificationGroup(grouped.notification, notification) {
				alreadyAdded = true
				grouped.count += 1
				break
			}
		}
		if !alreadyAdded {
			list = append(list, &groupedNotification{notification: notification, count: 1})
		}
	}
	return list
}

func isSameNotificationGroup(n1, n2 *models.Notification) bool {
	return n1.ResourceId == n2.ResourceId && n1.Reason == n2.Reason
}
