package printer

import "github.com/opf/openproject-cli/models"

// WhoamiEntry is one resolved identity (profile, server, user) for display.
type WhoamiEntry struct {
	Profile string
	Host    string
	User    *models.User
}

func Whoami(profile, host string, user *models.User) {
	activeRenderer.Whoami(profile, host, user)
}

// WhoamiList renders several identities as one output value, so JSON mode
// emits a single array instead of multiple top-level objects.
func WhoamiList(entries []WhoamiEntry) {
	activeRenderer.WhoamiList(entries)
}
