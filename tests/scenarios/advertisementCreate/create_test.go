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
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

const (
	testPhotoPath        = "testdata/photo.jpg"
	unsupportedPhotoPath = "testdata/unsupported.txt"
	invalidRequestError  = "invalid_request"
	unauthorizedError    = "unauthorized"
	incorrectToken       = "incorrect_token"
	negativePrice        = -100.0
	maxTitleLength       = 50
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
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
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
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Title = ""

	s.T().Log("Пытаемся создать объявление с пустым названием")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "title is required", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementNegativePrice() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Price = negativePrice

	s.T().Log("Пытаемся создать объявление с отрицательной ценой")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "price must be between 0 and 1000000 inclusive", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementLongTitle() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Title = strings.Repeat("A", maxTitleLength+1)

	s.T().Log("Пытаемся создать объявление с слишком длинным названием")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "title must not exceed 50 characters provided 51", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementWithoutToken() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)

	s.T().Log("Пытаемся создать объявление без токена")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), "", request, http.StatusUnauthorized,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), unauthorizedError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Authorization header required", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementWithIncorrectToken() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)

	s.T().Log("Пытаемся создать объявление с неверным токеном")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), incorrectToken, request, http.StatusUnauthorized,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), unauthorizedError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Invalid or expired token", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementWithoutPhotos() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Photos = nil

	s.T().Log("Пытаемся создать объявление без фотографий")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "at least 1 photo is required", resp.Message, "Некорректное сообщение об ошибке")
}

func (s *TestSuite) TestCreateAdvertisementEmptyDescription() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Description = ""

	s.T().Log("Пытаемся создать объявление с пустым описанием")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "description is required", resp.Message, "Некорректное сообщение об ошибке")
}

// TestCreateAdvertisementUnsupportedPhoto проверяет отказ при загрузке текстового файла вместо фотографии
func (s *TestSuite) TestCreateAdvertisementUnsupportedPhoto() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Photos = []string{unsupportedPhotoPath}

	s.T().Log("Пытаемся создать объявление с текстовым файлом вместо фотографии")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)
	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки неподдерживаемого формата фотографии")
	require.Equal(s.T(), `photo 1: file type not allowed: extension ".txt"`, resp.Message, "Некорректное сообщение об ошибке неподдерживаемого формата фотографии")
}

// TestCreateAdvertisementTooManyPhotos проверяет отказ при загрузке четырёх фотографий вместо разрешённых трёх
func (s *TestSuite) TestCreateAdvertisementTooManyPhotos() {
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	request.Photos = []string{
		testPhotoPath,
		testPhotoPath,
		testPhotoPath,
		testPhotoPath,
	}
	s.T().Log("Пытаемся создать объявление с четырьмя фотографиями при лимите три")
	resp := advertisements.CreateAdvertisementExpectError(
		s.T(), s.loginData.Token, request, http.StatusBadRequest,
	)
	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), invalidRequestError, resp.Error, "Некорректный код ошибки превышения количества фотографий")
	require.Equal(s.T(), "maximum 3 photos allowed", resp.Message, "Некорректное сообщение об ошибке превышения количества фотографий")
}
