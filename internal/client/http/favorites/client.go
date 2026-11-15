package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorites отправляет запрос добавления объявления в избранное
func HttpPostFavorites(t *testing.T, token string, request models.AddFavoriteRequest, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

// HttpGetFavorites отправляет запрос получения списка избранных объявлений
func HttpGetFavorites(t *testing.T, token string, expectedStatus, page int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Query("page", strconv.Itoa(page)).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного
func HttpDeleteFavorite(t *testing.T, token, ad_id string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + ad_id).
		Expect(t).
		Status(expectedStatus).
		End().Response
}
