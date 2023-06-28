package dtos

import (
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/models"
)

type ActivityDto struct {
	Id        uint64            `json:"id,omitempty"`
	Comment   *LongTextDto      `json:"comment,omitempty"`
	Details   []*LongTextDto    `json:"details,omitempty"`
	Version   uint64            `json:"version,omitempty"`
	CreatedAt string            `json:"createdAt,omitempty"`
	UpdatedAt string            `json:"updatedAt,omitempty"`
	Links     *activityLinksDto `json:"_links,omitempty"`
}

type activityElements struct {
	Elements []*ActivityDto `json:"elements"`
}

type ActivityCollectionDto struct {
	Embedded activityElements `json:"_embedded"`
	Type     string           `json:"_type"`
	Total    uint64           `json:"total"`
	Count    uint64           `json:"count"`
}

type activityLinksDto struct {
	User *activityUserLinkDto `json:"user"`
}

type activityUserLinkDto struct {
	Href string `json:"href"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *ActivityDto) Convert() (activity *models.Activity, err error) {
	var userId uint64
	if len(dto.Links.User.Href) > 0 {
		userHrefParts := strings.Split(dto.Links.User.Href, "/")
		userIdStr := userHrefParts[len(userHrefParts)-1]
		userId, err = strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &models.Activity{
		Id:        dto.Id,
		Comment:   dto.Comment.Raw,
		Details:   mapDetailsDto(dto.Details),
		Version:   dto.Version,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		UserId:    userId,
	}, nil
}

func (dto *ActivityCollectionDto) Convert() (act []*models.Activity, err error) {
	var activities = make([]*models.Activity, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		activities[idx], err = p.Convert()
		if err != nil {
			return nil, err
		}
	}

	return activities, nil
}

func mapDetailsDto(detailsDto []*LongTextDto) []*string {
	var details = make([]*string, len(detailsDto))
	for idx, d := range detailsDto {
		details[idx] = &d.Convert().Raw
	}

	return details
}
