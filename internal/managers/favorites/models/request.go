package models

// AddFavoriteRequest содержит id объявления для добавления в избранное
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}
