package printer

import (
	"fmt"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
)

func Users(users []*models.User) {
	var maxIdLength = 0
	for _, u := range users {
		maxIdLength = common.Max(maxIdLength, idLength(u.Id))
	}

	for _, u := range users {
		printUser(u, maxIdLength)
	}
}

func User(user *models.User) {
	printUser(user, idLength(user.Id))
}

func printUser(user *models.User, maxIdLength int) {
	var parts []string

	diff := maxIdLength - idLength(user.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), user.Id)
	parts = append(parts, Red(idStr))

	parts = append(parts, Cyan(user.Name))
	activePrinter.Println(strings.Join(parts, " "))
}
