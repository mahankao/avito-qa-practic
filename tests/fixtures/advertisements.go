package fixtures

import (
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

const (
	defaultTitlePrefix = "Тестовое объявление "
	defaultDescription = "Описание тестового объявления"
	defaultPrice       = 100.0
	defaultQuantity    = 10
	titleSuffixLength  = 8
)

// NewAdvertisementRequest возвращает новый корректный запрос создания объявления
func NewAdvertisementRequest(photoPath string) models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       defaultTitlePrefix + utils.RandomString(titleSuffixLength),
		Description: defaultDescription,
		Price:       defaultPrice,
		Quantity:    defaultQuantity,
		Photos:      []string{photoPath},
	}
}
