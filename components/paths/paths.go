package paths

import "fmt"

func Notifications() string {
	return Root() + "/notifications"
}

func Principals() string {
	return Root() + "/principals"
}

func Project(id string) string {
	return Projects() + "/" + id
}

func Projects() string {
	return Root() + "/projects"
}

func ProjectVersions(projectId string) string {
	return Project(projectId) + "/versions"
}

func ProjectWorkPackages(projectId string) string {
	return Project(projectId) + "/work_packages"
}

func Budget(id uint64) string {
	return Budgets() + fmt.Sprintf("/%d", id)
}

func Budgets() string {
	return Root() + "/budgets"
}

func ProjectBudgets(projectId string) string {
	return Project(projectId) + "/budgets"
}

func Root() string {
	return "/api/v3"
}

func Status() string {
	return Root() + "/statuses"
}

func TimeEntries() string {
	return Root() + "/time_entries"
}

func TimeEntry(id uint64) string {
	return TimeEntries() + fmt.Sprintf("/%d", id)
}

func TimeEntryActivities() string {
	return TimeEntries() + "/activities"
}

func Types() string {
	return Root() + "/types"
}

func User(id uint64) string {
	return Users() + fmt.Sprintf("/%d", id)
}

func UserMe() string {
	return Users() + "/me"
}

func Users() string {
	return Root() + "/users"
}

func WorkPackage(id uint64) string {
	return WorkPackages() + fmt.Sprintf("/%d", id)
}

func WorkPackages() string {
	return Root() + "/work_packages"
}

func WorkPackageActivities(id uint64) string {
	return WorkPackage(id) + "/activities"
}
