package advertisementCreate

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TestCreateAdvertisementPositive() {
	s.T().Log("Создаём объявление со всеми параметрами")
	request := adModels.CreateAdvertisementRequest{
		Title:       "Тестовая машина " + utils.RandomString(8),
		Description: "Описание тестовой машины",
		Price:       100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}
	ad := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request)
	defer advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)

	s.T().Log("Получаем объявление по ID и проверяем сохранённые параметры")
	savedAd := advertisements.GetAdvertisement(s.T(), s.loginData.Token, ad.ID)
	require.Equal(s.T(), request.Title, savedAd.Title, "Название объявления не совпадает")
	require.Equal(s.T(), request.Description, savedAd.Description, "Описание объявления не совпадает")
	require.Equal(s.T(), request.Price, savedAd.Price, "Цена объявления не совпадает")
	require.Equal(s.T(), request.Quantity, savedAd.Quantity, "Количество объявления не совпадает")

	s.T().Log("Проверяем поиск по полному названию")
	list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, request.Title)
	assertions.ContainsAdvertisementByTitle(s.T(), list.Items, request.Title)

	s.T().Log("Получаем фотографии и проверяем их данные")
	photos := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, ad.ID)
	require.Len(s.T(), photos, len(request.Photos), "Количество фотографий не совпадает")
	for _, photo := range photos {
		require.NotEmpty(s.T(), photo.ID, "У фотографии отсутствует ID")
		require.NotEmpty(s.T(), photo.URL, "У фотографии отсутствует URL")
	}
}

func (s *TestSuite) TestCreateAdvertisementEmptyTitle() {
	request := adModels.CreateAdvertisementRequest{
		Title:       "",
		Description: "Описание тестовой машины",
		Price:       100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}

	s.T().Log("Пытаемся создать объявление с пустым названием")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "invalid_request", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "title is required", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementNegativePrice() {
	request := adModels.CreateAdvertisementRequest{
		Title:       "Тестовая машина " + utils.RandomString(8),
		Description: "Описание тестовой машины",
		Price:       -100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}

	s.T().Log("Пытаемся создать объявление с отрицательной ценой")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "invalid_request", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "price must be between 0 and 1000000 inclusive", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementLongTitle() {
	request := adModels.CreateAdvertisementRequest{
		Title:       strings.Repeat("A", 51),
		Description: "Описание тестовой машины",
		Price:       100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}

	s.T().Log("Пытаемся создать объявление с слишком длинным названием")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "invalid_request", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "title must not exceed 50 characters provided 51", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementWithoutToken() {
	request := adModels.CreateAdvertisementRequest{
		Title:       "Тестовая машина " + utils.RandomString(8),
		Description: "Описание тестовой машины",
		Price:       100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}

	s.T().Log("Пытаемся создать объявление без токена")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), "", request, http.StatusUnauthorized,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "unauthorized", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Authorization header required", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementWithIncorrectToken() {
	request := adModels.CreateAdvertisementRequest{
		Title:       "Тестовая машина " + utils.RandomString(8),
		Description: "Описание тестовой машины",
		Price:       100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}

	s.T().Log("Пытаемся создать объявление с неверным токеном")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), "incorrect_token", request, http.StatusUnauthorized,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "unauthorized", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Invalid or expired token", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementWithoutPhotos() {
	request := adModels.CreateAdvertisementRequest{
		Title:       "Тестовая машина " + utils.RandomString(8),
		Description: "Описание тестовой машины",
		Price:       100.0,
		Quantity:    10,
	}

	s.T().Log("Пытаемся создать объявление без фотографий")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "invalid_request", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "at least 1 photo is required", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementEmptyDescription() {
	request := adModels.CreateAdvertisementRequest{
		Title:       "Тестовая машина " + utils.RandomString(8),
		Description: "",
		Price:       100.0,
		Quantity:    10,
		Photos:      []string{"testdata/photo.jpg"},
	}

	s.T().Log("Пытаемся создать объявление с пустым описанием")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "invalid_request", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "description is required", resp.Message, "Некорректное сообщение об ошибке")
}
