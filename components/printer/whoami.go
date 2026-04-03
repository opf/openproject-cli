package printer

import "github.com/opf/openproject-cli/models"

func Whoami(host string, user *models.User) {
	activeRenderer.Whoami(host, user)
}
