package assertions

import "api-tests-template/internal/managers/favorites/models"

// FavoriteContainsAdvertisement проверяет наличие объявления в списке избранного по ID
func FavoriteContainsAdvertisement(items []models.Favorite, ad_id string) bool {
	for _, item := range items {
		if item.Advertisement.ID == ad_id {
			return true
		}
	}
	return false
}
