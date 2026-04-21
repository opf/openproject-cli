package printer

import "github.com/opf/openproject-cli/models"

func Whoami(profile, host string, user *models.User) {
	activeRenderer.Whoami(profile, host, user)
}
