package users

import (
	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func ByIds(ids []uint64) []*models.User {
	if len(ids) == 0 {
		return []*models.User{}
	}
	var filters []requests.Filter
	filters = append(filters, IdFilter(ids))

	query := requests.NewFilterQuery(filters)

	requestUrl := paths.Principals()

	response, err := requests.Get(requestUrl, &query)
	if err != nil {
		printer.Error(err)
	}

	userCollection := parser.Parse[dtos.UserCollectionDto](response)
	return userCollection.Convert()
}

func Me() (*models.User, error) {
	response, err := requests.Get(paths.UserMe(), nil)
	if err != nil {
		return nil, err
	}

	user := parser.Parse[dtos.UserDto](response)
	return user.Convert(), nil
}

func LookUp(id uint64) (*models.User, error) {
	response, err := requests.Get(paths.User(id), nil)
	if err != nil {
		return nil, err
	}

	user := parser.Parse[dtos.UserDto](response)
	return user.Convert(), nil
}

func Search(input string) ([]*models.User, error) {
	inputAsId, userId := common.ParseId(input)
	var filters []requests.Filter

	if inputAsId {
		filters = append(filters, IdFilter([]uint64{userId}))
	} else {
		filters = append(filters, resources.TypeAheadFilter(input))
	}

	query := requests.NewFilterQuery(filters)

	response, err := requests.Get(paths.Principals(), &query)
	if err != nil {
		return nil, err
	}

	userCollection := parser.Parse[dtos.UserCollectionDto](response)
	return userCollection.Convert(), nil
}
