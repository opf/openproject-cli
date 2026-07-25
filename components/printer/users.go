package printer

import "github.com/opf/openproject-cli/models"

func Users(users []*models.User) {
	activeRenderer.Users(users)
}

func User(user *models.User) {
	activeRenderer.User(user)
}
