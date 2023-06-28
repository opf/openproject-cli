package dtos

import "github.com/opf/openproject-cli/models"

type UserDto struct {
	Id        uint64 `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type userElements struct {
	Elements []*UserDto `json:"elements"`
}

type UserCollectionDto struct {
	Embedded userElements `json:"_embedded"`
	Type     string       `json:"_type"`
	Total    uint64       `json:"total"`
	Count    uint64       `json:"count"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *UserDto) Convert() *models.User {
	return &models.User{
		Id:        dto.Id,
		Name:      dto.Name,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
}

func (dto *UserCollectionDto) Convert() []*models.User {
	var users = make([]*models.User, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		users[idx] = p.Convert()
	}

	return users
}
