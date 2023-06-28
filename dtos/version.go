package dtos

import (
	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
)

type VersionDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type versionElements struct {
	Elements []*VersionDto `json:"elements"`
}

type VersionCollectionDto struct {
	Embedded *versionElements `json:"_embedded"`
	Type     string           `json:"_type"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *VersionDto) Convert() *models.Version {
	return &models.Version{
		Id:   uint64(dto.Id),
		Name: dto.Name,
	}
}

func (dto *VersionCollectionDto) Convert() []*models.Version {
	return common.Reduce(
		dto.Embedded.Elements,
		func(state []*models.Version, version *VersionDto) []*models.Version {
			return append(state, version.Convert())
		},
		[]*models.Version{},
	)
}
