package models

import adModels "api-tests-template/internal/managers/advertisements/models"

// Favorite описывает объявление, добавленное в избранное
type Favorite struct {
	FavoriteID    string                 `json:"favorite_id"`
	AddedAt       string                 `json:"added_at"`
	Advertisement adModels.Advertisement `json:"advertisement"`
}

// FavoritesList описывает список избранного
type FavoritesList struct {
	Items      []Favorite `json:"items"`
	TotalPages int        `json:"total_pages"`
}
