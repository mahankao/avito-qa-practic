package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное и возвращает запись
func AddFavorite(t *testing.T, token string, request models.AddFavoriteRequest) models.Favorite {
	raw := favorites.HttpPostFavorites(t, token, request, http.StatusCreated)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Добавление в избранное: не удалось прочитать ответ")

	resp := models.Favorite{}
	require.NoError(t, json.Unmarshal(body, &resp), "Добавление в избранное: не удалось разобрать JSON")

	return resp
}

// GetFavoritesList возвращает избранные объявления со всех страниц
func GetFavoritesList(t *testing.T, token string) models.FavoritesList {
	result := models.FavoritesList{}
	for page := 1; ; page++ {
		raw := favorites.HttpGetFavorites(t, token, http.StatusOK, page)
		body, err := io.ReadAll(raw.Body)
		closeErr := raw.Body.Close()
		require.NoError(t, err, "Получение избранного: не удалось прочитать ответ")
		require.NoError(t, closeErr, "Получение избранного: не удалось закрыть тело ответа")

		resp := models.FavoritesList{}
		require.NoError(t, json.Unmarshal(body, &resp), "Получение избранного: не удалось разобрать JSON")
		result.Items = append(result.Items, resp.Items...)
		result.TotalPages = resp.TotalPages
		if page >= resp.TotalPages {
			return result
		}
	}
}

// DeleteFavorite удаляет объявление из избранного
func DeleteFavorite(t *testing.T, token, ad_id string) {
	raw := favorites.HttpDeleteFavorite(t, token, ad_id, http.StatusNoContent)
	defer raw.Body.Close()
}

// GetFavoritesExpectError возвращает ошибку получения списка избранного
func GetFavoritesExpectError(t *testing.T, token string, expectedStatus int) adModels.ErrorResponse {
	raw := favorites.HttpGetFavorites(t, token, expectedStatus, 1)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Ошибка получения избранного: не удалось прочитать ответ")

	resp := adModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка получения избранного: не удалось разобрать JSON")

	return resp
}

// AddFavoriteExpectError возвращает ошибку добавления в избранное
func AddFavoriteExpectError(t *testing.T, token string, request models.AddFavoriteRequest, expectedStatus int) adModels.ErrorResponse {
	raw := favorites.HttpPostFavorites(t, token, request, expectedStatus)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Ошибка добавления в избранное: не удалось прочитать ответ")

	resp := adModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка добавления в избранное: не удалось разобрать JSON")

	return resp
}

// DeleteFavoriteExpectError возвращает ошибку удаления из избранного
func DeleteFavoriteExpectError(t *testing.T, token, ad_id string, expectedStatus int) adModels.ErrorResponse {
	raw := favorites.HttpDeleteFavorite(t, token, ad_id, expectedStatus)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Ошибка удаления из избранного: не удалось прочитать ответ")

	resp := adModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка удаления из избранного: не удалось разобрать JSON")

	return resp
}
