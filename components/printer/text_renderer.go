package printer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
)

type TextRenderer struct{}

func (r *TextRenderer) Budget(b *models.Budget) {
	printBudget(b, idLength(b.Id))
}

func (r *TextRenderer) Budgets(bs []*models.Budget) {
	var maxIdLength = 0
	for _, b := range bs {
		maxIdLength = common.Max(maxIdLength, idLength(b.Id))
	}
	for _, b := range bs {
		printBudget(b, maxIdLength)
	}
}

func (r *TextRenderer) WorkPackage(wp *models.WorkPackage) {
	printHeadline(wp, idLength(wp.Id), 0, utf8.RuneCountInString(wp.Type))
	printAttributes(wp)
	activePrinter.Println()
	printOpenLink(wp)
	activePrinter.Println()
	printDescription(wp)
}

func (r *TextRenderer) WorkPackages(wps []*models.WorkPackage) {
	var maxIdLength = 0
	var maxTypeLength = 0
	var maxStatusLength = 0
	for _, w := range wps {
		maxIdLength = common.Max(maxIdLength, idLength(w.Id))
		maxTypeLength = common.Max(maxTypeLength, utf8.RuneCountInString(w.Type))
		maxStatusLength = common.Max(maxStatusLength, utf8.RuneCountInString(w.Status))
	}
	for _, wp := range wps {
		printHeadline(wp, maxIdLength, maxStatusLength, maxTypeLength)
	}
}

func (r *TextRenderer) Project(p *models.Project) {
	printProject(p)
}

func (r *TextRenderer) Projects(ps []*models.Project) {
	for _, p := range ps {
		printProject(p)
	}
}

func (r *TextRenderer) User(u *models.User) {
	printUser(u, idLength(u.Id))
}

func (r *TextRenderer) Users(us []*models.User) {
	var maxIdLength = 0
	for _, u := range us {
		maxIdLength = common.Max(maxIdLength, idLength(u.Id))
	}
	for _, u := range us {
		printUser(u, maxIdLength)
	}
}

func (r *TextRenderer) Types(types []*models.Type) {
	var maxIdLength = 0
	for _, t := range types {
		maxIdLength = common.Max(maxIdLength, idLength(t.Id))
	}
	for _, t := range types {
		printType(t, maxIdLength)
	}
}

func (r *TextRenderer) Status(s *models.Status) {
	printStatus(s, idLength(s.Id))
}

func (r *TextRenderer) StatusList(statuses []*models.Status) {
	var maxIdLength = 0
	for _, s := range statuses {
		maxIdLength = common.Max(maxIdLength, idLength(s.Id))
	}
	for _, s := range statuses {
		printStatus(s, maxIdLength)
	}
}

func (r *TextRenderer) TimeEntry(t *models.TimeEntry) {
	printTimeEntry(t, idLength(t.Id), len(t.Activity), len(t.Project))
}

func (r *TextRenderer) TimeEntryList(entries []*models.TimeEntry) {
	var maxIdLength = 0
	var maxActivityLength = 0
	var maxProjectLength = 0
	for _, t := range entries {
		maxIdLength = common.Max(maxIdLength, idLength(t.Id))
		maxActivityLength = common.Max(maxActivityLength, len(t.Activity))
		maxProjectLength = common.Max(maxProjectLength, len(t.Project))
	}
	for _, t := range entries {
		printTimeEntry(t, maxIdLength, maxActivityLength, maxProjectLength)
	}
}

func (r *TextRenderer) Notification(n *models.Notification) {
	printNotification(&groupedNotification{notification: n, count: 1}, idLength(n.ResourceId), len(n.Reason))
}

func (r *TextRenderer) Notifications(ns []*models.Notification) {
	grouped := group(ns)
	var maxIdLength int
	var maxReasonLength int
	for _, element := range grouped {
		maxIdLength = common.Max(maxIdLength, idLength(element.notification.ResourceId))
		maxReasonLength = common.Max(maxReasonLength, len(element.notification.Reason))
	}
	for _, n := range grouped {
		printNotification(n, maxIdLength, maxReasonLength)
	}
}

func (r *TextRenderer) Activities(activities []*models.Activity, users []*models.User) {
	for _, activity := range activities {
		user := &models.User{}
		if activity.UserId > 0 {
			idx := sort.Search(len(users), func(i int) bool { return users[i].Id >= activity.UserId })
			if idx < len(users) && users[idx].Id == activity.UserId {
				user = users[idx]
			}
		}
		printActivityHeadline(activity, user)
		printActivityBody(activity)
		activePrinter.Println("")
	}
}

func (r *TextRenderer) CustomActions(actions []*models.CustomAction) {
	for _, a := range actions {
		printCustomAction(a)
	}
}

func (r *TextRenderer) Number(n int64) {
	activePrinter.Printf("%s\n", Cyan(strconv.FormatInt(n, 10)))
}

func printBudget(b *models.Budget, maxIdLength int) {
	diff := maxIdLength - idLength(b.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), b.Id)
	activePrinter.Printf("%s %s\n", Red(idStr), Cyan(b.Subject))
}

func printProject(p *models.Project) {
	id := fmt.Sprintf("#%d", p.Id)
	activePrinter.Printf("%s %s\n", Red(id), Cyan(p.Name))
}

func printUser(u *models.User, maxIdLength int) {
	diff := maxIdLength - idLength(u.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), u.Id)
	activePrinter.Println(strings.Join([]string{Red(idStr), Cyan(u.Name)}, " "))
}

func printType(t *models.Type, maxIdLength int) {
	diff := maxIdLength - idLength(t.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), t.Id)
	activePrinter.Println(strings.Join([]string{Red(idStr), Cyan(t.Name)}, " "))
}

func printStatus(s *models.Status, maxIdLength int) {
	diff := maxIdLength - idLength(s.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), s.Id)
	parts := []string{Red(idStr), Cyan(s.Name)}
	if s.IsDefault {
		parts = append(parts, fmt.Sprintf("(%s)", Yellow("default")))
	}
	activePrinter.Println(strings.Join(parts, " "))
}

func printCustomAction(a *models.CustomAction) {
	activePrinter.Printf("%s %s\n", Red(fmt.Sprintf("#%d", a.Id)), Cyan(a.Name))
}
