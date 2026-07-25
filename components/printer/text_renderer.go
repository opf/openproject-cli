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
	printHeadline(wp, displayIdLength(wp), 0, utf8.RuneCountInString(wp.Type))
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
		maxIdLength = common.Max(maxIdLength, displayIdLength(w))
		maxTypeLength = common.Max(maxTypeLength, utf8.RuneCountInString(w.Type))
		maxStatusLength = common.Max(maxStatusLength, utf8.RuneCountInString(w.Status))
	}
	for _, wp := range wps {
		printHeadline(wp, maxIdLength, maxStatusLength, maxTypeLength)
	}
}

func (r *TextRenderer) WorkPackageDetails(p *models.WorkPackageInspectPayload) {
	wp := p.WorkPackage
	idStr := fmt.Sprintf("#%d", wp.ID)
	activePrinter.Printf("%s %s %s\n", Red(idStr), Green(strings.ToUpper(wp.Type)), Cyan(wp.Subject))
	activePrinter.Printf("[%s]\n", Yellow(wp.Status))

	assigneeStr := wp.Assignee
	if len(assigneeStr) == 0 {
		assigneeStr = "-"
	}
	activePrinter.Printf("Assignee: %s\n", assigneeStr)
	activePrinter.Printf("Project: %s (%s)\n", wp.Project.Name, wp.Project.Identifier)
	if wp.ParentID != nil {
		activePrinter.Printf("Parent: #%d\n", *wp.ParentID)
	}

	printCustomFields(wp.FieldLabels, wp.Fields)

	activePrinter.Printf("Children (%d):\n", len(p.Children))
	for _, child := range p.Children {
		activePrinter.Printf("#%d [%s] %s\n", child.ID, child.Status, child.Subject)
	}
}

func (r *TextRenderer) WorkPackageCreatePlan(p *models.WorkPackageCreatePlan) {
	activePrinter.Println("Dry run — no changes applied.")
	printPlanField("operation", p.Operation)
	printPlanField("project_id", p.ProjectID)
	if p.ParentID != nil {
		printPlanField("parent_id", strconv.FormatUint(*p.ParentID, 10))
	}
	printPlanField("subject", p.WorkPackage.Subject)
	printPlanField("type", p.WorkPackage.Type)
	printPlanField("description", p.WorkPackage.Description)
}

func (r *TextRenderer) WorkPackageUpdatePlan(p *models.WorkPackageUpdatePlan) {
	activePrinter.Println("Dry run — no changes applied.")
	printPlanField("operation", p.Operation)
	printPlanField("work_package_id", p.WorkPackageID)
	printPlanField("subject", p.Subject)
	printPlanField("type", p.Type)
	printPlanField("assignee", p.Assignee)
	printPlanField("status", p.Status)
	if p.Description != nil {
		printPlanField("description", *p.Description)
	}
	printPlanField("action", p.Action)
	printPlanField("attach", p.Attach)
	printResolvedFields(p.ResolvedFields)
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

func (r *TextRenderer) Whoami(profile, host string, user *models.User) {
	activePrinter.Printf("Profile: %s\n", Yellow(profile))
	activePrinter.Printf("Server:  %s\n", Cyan(host))
	activePrinter.Printf("User:    %s %s\n", Red(fmt.Sprintf("#%d", user.Id)), Cyan(user.Name))
}

func (r *TextRenderer) WhoamiList(entries []WhoamiEntry) {
	for i, entry := range entries {
		if i > 0 {
			activePrinter.Println()
		}
		r.Whoami(entry.Profile, entry.Host, entry.User)
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
	activePrinter.Printf("%s %s (%s)\n", Red(id), Cyan(p.Name), p.Identifier)
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

// printCustomFields prints one "label: value" line per API field, sorted
// first by label and then by API name so output stays deterministic even
// though a single label can map to more than one API field.
func printCustomFields(fieldLabels map[string][]string, fields map[string]any) {
	labels := make([]string, 0, len(fieldLabels))
	for label := range fieldLabels {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	for _, label := range labels {
		apiNames := append([]string(nil), fieldLabels[label]...)
		sort.Strings(apiNames)
		for _, apiName := range apiNames {
			activePrinter.Printf("%s: %v\n", label, fields[apiName])
		}
	}
}

// printPlanField prints a single "name: value" line, skipping fields that
// are empty or omitted from the dry-run plan.
func printPlanField(name, value string) {
	if len(value) == 0 {
		return
	}
	activePrinter.Printf("%s: %s\n", name, value)
}

// printResolvedFields prints one line per schema-resolved custom field
// assignment, sorted by key for deterministic output.
func printResolvedFields(resolvedFields map[string]models.ResolvedField) {
	keys := make([]string, 0, len(resolvedFields))
	for key := range resolvedFields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		field := resolvedFields[key]
		activePrinter.Printf("%s (%s): %v\n", key, field.APIField, field.Value)
	}
}
