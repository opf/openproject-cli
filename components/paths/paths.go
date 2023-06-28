package paths

import "fmt"

func Notifications() string {
	return Root() + "/notifications"
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
	return "api/v3"
}

func WorkPackage(id uint64) string {
	return WorkPackages() + fmt.Sprintf("/%d", id)
}

func WorkPackages() string {
	return Root() + "/work_packages"
}
