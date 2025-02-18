package paths

import "fmt"

func Notifications() string {
	return Root() + "/notifications"
}

func Principals() string {
	return Root() + "/principals"
}

func Project(id uint64) string {
	return Projects() + fmt.Sprintf("/%d", id)
}

func Projects() string {
	return Root() + "/projects"
}

func ProjectVersions(projectId uint64) string {
	return Project(projectId) + "/versions"
}

func ProjectWorkPackages(projectId uint64) string {
	return Project(projectId) + "/work_packages"
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
