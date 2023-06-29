package dtos

import (
	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
)

type typeLinks struct {
	Self *LinkDto `json:"self"`
}

type TypeDto struct {
	Id    int64      `json:"id"`
	Name  string     `json:"name"`
	Links *typeLinks `json:"_links"`
}

type typeElements struct {
	Elements []*TypeDto `json:"elements"`
}

type TypeCollectionDto struct {
	Embedded *typeElements `json:"_embedded"`
	Type     string        `json:"_type"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *TypeDto) Convert() *models.Type {
	return &models.Type{
		Id:   uint64(dto.Id),
		Name: dto.Name,
	}
}

func (dto *TypeCollectionDto) Convert() []*models.Type {
	var list = make([]*models.Type, len(dto.Embedded.Elements))

	for idx, element := range dto.Embedded.Elements {
		list[idx] = element.Convert()
	}

	return list
}

type TypeDtos []*TypeDto

func (dtos TypeDtos) Convert() []*models.Type {
	return common.Reduce(dtos,
		func(acc []*models.Type, dto *TypeDto) []*models.Type {
			return append(acc, dto.Convert())
		}, []*models.Type{})
}
