package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// GetAdvertisementsList возвращает объявления по поисковой строке
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Поиск объявлений: не удалось прочитать ответ")

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp), "Поиск объявлений: не удалось разобрать JSON")

	return resp
}

// CreateAdvertisement создаёт объявление и возвращает его данные
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	raw := advertisements.HttpPostAdvertisement(t, token, request, http.StatusCreated)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Создание объявления: не удалось прочитать ответ")

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp), "Создание объявления: не удалось разобрать JSON")

	return resp
}

// CreateAdvertisementExpectError отправляет запрос создания и возвращает ошибку API
func CreateAdvertisementExpectError(t *testing.T, token string, request models.CreateAdvertisementRequest, expectedStatus int) models.ErrorResponse {
	raw := advertisements.HttpPostAdvertisement(t, token, request, expectedStatus)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Ошибка создания объявления: не удалось прочитать ответ")

	resp := models.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка создания объявления: не удалось разобрать JSON")

	return resp
}

// GetAdvertisement возвращает данные объявления по ID
func GetAdvertisement(t *testing.T, token string, id string) models.Advertisement {
	raw := advertisements.HttpGetAdvertisement(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Получение объявления: не удалось прочитать ответ")

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp), "Получение объявления: не удалось разобрать JSON")

	return resp
}

// DeleteAdvertisement удаляет объявление по ID
func DeleteAdvertisement(t *testing.T, token string, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()
}

// GetAdvertisementPhotos возвращает фотографии объявления по ID
func GetAdvertisementPhotos(t *testing.T, token string, id string) []models.Photo {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Получение фотографий: не удалось прочитать ответ")

	resp := []models.Photo{}
	require.NoError(t, json.Unmarshal(body, &resp), "Получение фотографий: не удалось разобрать JSON")

	return resp
}
