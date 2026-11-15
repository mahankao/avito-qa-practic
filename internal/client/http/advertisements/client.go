package advertisements

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements ищет объявления по подстроке названия
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement отправляет запрос создания объявления и проверяет ожидаемый статус
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description).
		MultipartFormData("price", fmt.Sprint(request.Price)).
		MultipartFormData("quantity", fmt.Sprint(request.Quantity)).
		MultipartFile("photos", request.Photos...).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

// HttpGetAdvertisement получает объявление по id
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по id
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusNoContent).
		End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления по id
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath + "/" + id + "/photos").
		Expect(t).
		Status(http.StatusOK).
		End().Response
}
