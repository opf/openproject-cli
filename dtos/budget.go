package dtos

import "github.com/opf/openproject-cli/models"

type BudgetDto struct {
	Id      int64  `json:"id"`
	Subject string `json:"subject"`
}

type budgetElements struct {
	Elements []*BudgetDto `json:"elements"`
}

type BudgetCollectionDto struct {
	Embedded *budgetElements `json:"_embedded"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *BudgetCollectionDto) Convert() []*models.Budget {
	if dto.Embedded == nil {
		return []*models.Budget{}
	}
	budgets := make([]*models.Budget, len(dto.Embedded.Elements))
	for i, b := range dto.Embedded.Elements {
		budgets[i] = b.Convert()
	}
	return budgets
}

func (dto *BudgetDto) Convert() *models.Budget {
	return &models.Budget{
		Id:      uint64(dto.Id),
		Subject: dto.Subject,
	}
}
