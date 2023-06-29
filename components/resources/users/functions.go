package users

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

const apiPath = "api/v3"
const usersPath = apiPath + "/principals"

func ByIds(ids []uint64) []*models.User {
	if len(ids) == 0 {
		return []*models.User{}
	}
	var filters []requests.Filter
	filters = append(filters, IdFilter(ids))

	query := requests.NewQuery(filters)

	requestUrl := usersPath

	response, err := requests.Get(requestUrl, &query)
	if err != nil {
		printer.Error(err)
	}

	userCollection := parser.Parse[dtos.UserCollectionDto](response)
	return userCollection.Convert()
}
