package printer

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/opf/openproject-cli/models"
)

type JsonRenderer struct{}

func (r *JsonRenderer) WorkPackage(wp *models.WorkPackage) {
	printJson(struct {
		Id          uint64 `json:"id"`
		Subject     string `json:"subject"`
		Type        string `json:"type"`
		Status      string `json:"status"`
		Assignee    string `json:"assignee"`
		Description string `json:"description"`
	}{wp.Id, wp.Subject, wp.Type, wp.Status, wp.Assignee, wp.Description})
}

func (r *JsonRenderer) WorkPackages(wps []*models.WorkPackage) {
	type item struct {
		Id       uint64 `json:"id"`
		Subject  string `json:"subject"`
		Type     string `json:"type"`
		Status   string `json:"status"`
		Assignee string `json:"assignee"`
	}
	out := make([]item, len(wps))
	for i, wp := range wps {
		out[i] = item{wp.Id, wp.Subject, wp.Type, wp.Status, wp.Assignee}
	}
	printJson(out)
}

func (r *JsonRenderer) Project(p *models.Project) {
	printJson(struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}{p.Id, p.Name})
}

func (r *JsonRenderer) Projects(ps []*models.Project) {
	type item struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
	out := make([]item, len(ps))
	for i, p := range ps {
		out[i] = item{p.Id, p.Name}
	}
	printJson(out)
}

func (r *JsonRenderer) User(u *models.User) {
	printJson(struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{u.Id, u.Name, u.FirstName, u.LastName})
}

func (r *JsonRenderer) Users(us []*models.User) {
	type item struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	out := make([]item, len(us))
	for i, u := range us {
		out[i] = item{u.Id, u.Name, u.FirstName, u.LastName}
	}
	printJson(out)
}

func (r *JsonRenderer) Types(types []*models.Type) {
	type item struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
	out := make([]item, len(types))
	for i, t := range types {
		out[i] = item{t.Id, t.Name}
	}
	printJson(out)
}

func (r *JsonRenderer) Status(s *models.Status) {
	printJson(struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		IsDefault bool   `json:"is_default"`
		IsClosed  bool   `json:"is_closed"`
	}{s.Id, s.Name, s.IsDefault, s.IsClosed})
}

func (r *JsonRenderer) StatusList(statuses []*models.Status) {
	type item struct {
		Id        uint64 `json:"id"`
		Name      string `json:"name"`
		IsDefault bool   `json:"is_default"`
		IsClosed  bool   `json:"is_closed"`
	}
	out := make([]item, len(statuses))
	for i, s := range statuses {
		out[i] = item{s.Id, s.Name, s.IsDefault, s.IsClosed}
	}
	printJson(out)
}

func (r *JsonRenderer) TimeEntry(t *models.TimeEntry) {
	printJson(struct {
		Id          uint64  `json:"id"`
		Comment     string  `json:"comment"`
		Project     string  `json:"project"`
		WorkPackage string  `json:"work_package"`
		SpentOn     string  `json:"spent_on"`
		Hours       float64 `json:"hours"`
		Activity    string  `json:"activity"`
		User        string  `json:"user"`
	}{t.Id, t.Comment, t.Project, t.WorkPackage, t.SpentOn.Format("2006-01-02"), t.Hours.Hours(), t.Activity, t.User})
}

func (r *JsonRenderer) TimeEntryList(entries []*models.TimeEntry) {
	type item struct {
		Id          uint64  `json:"id"`
		Comment     string  `json:"comment"`
		Project     string  `json:"project"`
		WorkPackage string  `json:"work_package"`
		SpentOn     string  `json:"spent_on"`
		Hours       float64 `json:"hours"`
		Activity    string  `json:"activity"`
		User        string  `json:"user"`
	}
	out := make([]item, len(entries))
	for i, t := range entries {
		out[i] = item{
			t.Id,
			t.Comment,
			t.Project,
			t.WorkPackage,
			t.SpentOn.Format("2006-01-02"),
			t.Hours.Hours(),
			t.Activity,
			t.User,
		}
	}
	printJson(out)
}

func (r *JsonRenderer) Notification(n *models.Notification) {
	printJson(struct {
		Id              uint64 `json:"id"`
		ResourceId      uint64 `json:"resource_id"`
		ResourceSubject string `json:"resource_subject"`
		Reason          string `json:"reason"`
		Read            bool   `json:"read"`
	}{n.Id, n.ResourceId, n.ResourceSubject, n.Reason, n.Read})
}

func (r *JsonRenderer) Notifications(ns []*models.Notification) {
	type item struct {
		Id              uint64 `json:"id"`
		ResourceId      uint64 `json:"resource_id"`
		ResourceSubject string `json:"resource_subject"`
		Reason          string `json:"reason"`
		Read            bool   `json:"read"`
	}
	out := make([]item, len(ns))
	for i, n := range ns {
		out[i] = item{n.Id, n.ResourceId, n.ResourceSubject, n.Reason, n.Read}
	}
	printJson(out)
}

func (r *JsonRenderer) Activities(activities []*models.Activity, users []*models.User) {
	type item struct {
		Id        uint64    `json:"id"`
		Comment   string    `json:"comment"`
		Details   []*string `json:"details"`
		UpdatedAt string    `json:"updated_at"`
		User      string    `json:"user"`
	}
	out := make([]item, len(activities))
	for i, a := range activities {
		userName := ""
		if a.UserId > 0 {
			idx := sort.Search(len(users), func(j int) bool { return users[j].Id >= a.UserId })
			if idx < len(users) && users[idx].Id == a.UserId {
				userName = users[idx].Name
			}
		}
		out[i] = item{a.Id, a.Comment, a.Details, a.UpdatedAt, userName}
	}
	printJson(out)
}

func (r *JsonRenderer) CustomActions(actions []*models.CustomAction) {
	type item struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
	}
	out := make([]item, len(actions))
	for i, a := range actions {
		out[i] = item{a.Id, a.Name}
	}
	printJson(out)
}

func (r *JsonRenderer) Number(n int64) {
	printJson(struct {
		Total int64 `json:"total"`
	}{n})
}

func printJson(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("{\"error\": \"failed to serialize output: %s\"}\n", err)
		return
	}
	fmt.Println(string(b))
}
