package configuration

import "time"

type Version struct {
	Version string
	Commit  string
	Date    time.Time
}

var CliVersion *Version

func Init(version *Version) {
	CliVersion = version
}

func BuildCliVersion(version string, commit string, date string) *Version {
	buildCommit := string([]rune(commit)[:7])

	var buildDate time.Time
	if date == "unknown" {
		buildDate = time.Now()
	} else {
		buildDate, _ = time.Parse(time.RFC3339, date)
	}

	return &Version{
		Version: version,
		Commit:  buildCommit,
		Date:    buildDate,
	}
}
