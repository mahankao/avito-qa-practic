package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	favModels "api-tests-template/internal/managers/favorites/models"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

const (
	testPhotoPath               = "../advertisementCreate/testdata/photo.jpg"
	nonexistentAdvertisementID  = "12345678-1234-4123-8123-123456789012"
	unauthorizedError           = "unauthorized"
	missingAuthorizationMessage = "Authorization header required"
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
	base.Precondition("Авторизуемся для тестов избранного")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TestFavoritesPositive() {
	s.T().Log("Подготовка: авторизуем второго пользователя")
	secondUser := auth.Login(s.T(), os.Getenv("TEST_SECOND_LOGIN"), os.Getenv("TEST_SECOND_PASSWORD"))
	require.NotEqual(s.T(), s.loginData.User.Id, secondUser.User.Id, "Для теста нужны два разных аккаунта")

	s.T().Log("Подготовка: первый пользователь создаёт объявление")
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	ad := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request)
	adDeleted := false
	defer func() {
		if !adDeleted {
			advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)
		}
	}()
	favoriteAdded := false
	defer func() {
		if favoriteAdded {
			favorites.DeleteFavorite(s.T(), secondUser.Token, ad.ID)
		}
	}()
	favoriteRequest := favModels.AddFavoriteRequest{AdvertisementID: ad.ID}

	s.T().Log("Добавляем объявление в избранное второго пользователя")
	favorite := favorites.AddFavorite(s.T(), secondUser.Token, favoriteRequest)
	favoriteAdded = true
	require.Equal(s.T(), ad.ID, favorite.Advertisement.ID, "В избранное добавилось другое объявление")
	require.NotEmpty(s.T(), favorite.FavoriteID, "Отсутствует ID записи избранного")

	list := favorites.GetFavoritesList(s.T(), secondUser.Token)
	require.True(s.T(), assertions.FavoriteContainsAdvertisement(list.Items, ad.ID), "Добавленное объявление отсутствует в избранном")
	s.T().Log("Проверено: добавленное объявление есть в списке")

	s.T().Log("Удаляем объявление из избранного")
	favorites.DeleteFavorite(s.T(), secondUser.Token, ad.ID)
	favoriteAdded = false
	list = favorites.GetFavoritesList(s.T(), secondUser.Token)
	require.False(s.T(), assertions.FavoriteContainsAdvertisement(list.Items, ad.ID), "Удалённая запись осталась в избранном")
	s.T().Log("Проверено: после удаления из избранного записи нет")

	s.T().Log("Повторно добавляем объявление и проверяем его наличие в избранном")
	favorites.AddFavorite(s.T(), secondUser.Token, favoriteRequest)
	favoriteAdded = true
	list = favorites.GetFavoritesList(s.T(), secondUser.Token)
	require.True(s.T(), assertions.FavoriteContainsAdvertisement(list.Items, ad.ID),
		"Перед удалением объявления его запись отсутствует в избранном")
	s.T().Log("Проверено: после повторного добавления запись вернулась в избранное")

	s.T().Log("Удаляем само объявление от имени владельца")
	advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)
	adDeleted = true
	list = favorites.GetFavoritesList(s.T(), secondUser.Token)
	favoriteAdded = assertions.FavoriteContainsAdvertisement(list.Items, ad.ID)
	require.False(s.T(), favoriteAdded, "После удаления самого объявления запись осталась в избранном")
	s.T().Log("Проверено: удалённое объявление исчезло из избранного")
}

func (s *TestSuite) TestAddFavoriteNonexistentAdvertisement() {
	request := favModels.AddFavoriteRequest{
		AdvertisementID: nonexistentAdvertisementID,
	}

	s.T().Log("Пытаемся добавить несуществующее объявление в избранное")
	resp := favorites.AddFavoriteExpectError(
		s.T(), s.loginData.Token, request, http.StatusNotFound,
	)

	s.T().Logf("Ответ API: error=%q, message=%q, details=%q",
		resp.Error, resp.Message, resp.Details)

	require.Equal(s.T(), "not_found", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Advertisement not found", resp.Message, "Некорректное сообщение об ошибке")
	s.T().Log("Проверены статус 404 и ошибка отсутствующего объявления")
}

func (s *TestSuite) TestDeleteFavoriteNonexistentAdvertisement() {
	s.T().Log("Пытаемся удалить несуществующее объявление из избранного")
	resp := favorites.DeleteFavoriteExpectError(
		s.T(), s.loginData.Token, nonexistentAdvertisementID, http.StatusNotFound,
	)
	s.T().Logf("Ответ API: error=%q, message=%q", resp.Error, resp.Message)

	require.Equal(s.T(), "not_found", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Favorite entry not found", resp.Message, "Некорректное сообщение об ошибке")
	s.T().Log("Проверены статус 404 и ошибка отсутствующей записи избранного")
}

func (s *TestSuite) TestAddFavoriteWithoutToken() {
	request := favModels.AddFavoriteRequest{AdvertisementID: nonexistentAdvertisementID}
	s.T().Log("Пытаемся добавить объявление в избранное без авторизации")
	resp := favorites.AddFavoriteExpectError(s.T(), "", request, http.StatusUnauthorized)

	require.Equal(s.T(), unauthorizedError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), missingAuthorizationMessage, resp.Message, "Некорректное сообщение об ошибке")
	s.T().Log("POST без токена вернул ожидаемые статус 401 и сообщение")
}

func (s *TestSuite) TestGetFavoritesWithoutToken() {
	s.T().Log("Пытаемся получить избранное без авторизации")
	resp := favorites.GetFavoritesExpectError(s.T(), "", http.StatusUnauthorized)

	require.Equal(s.T(), unauthorizedError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), missingAuthorizationMessage, resp.Message, "Некорректное сообщение об ошибке")
	s.T().Log("GET без токена вернул ожидаемые статус 401 и сообщение")
}

func (s *TestSuite) TestDeleteFavoriteWithoutToken() {
	s.T().Log("Пытаемся удалить объявление из избранного без авторизации")
	resp := favorites.DeleteFavoriteExpectError(s.T(), "", nonexistentAdvertisementID, http.StatusUnauthorized)

	require.Equal(s.T(), unauthorizedError, resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), missingAuthorizationMessage, resp.Message, "Некорректное сообщение об ошибке")
	s.T().Log("DELETE без токена вернул ожидаемые статус 401 и сообщение")
}

func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	s.T().Log("Создаём собственное объявление для проверки запрета добавления в избранное")
	request := fixtures.NewAdvertisementRequest(testPhotoPath)
	ad := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request)
	defer advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)

	s.T().Log("Пытаемся добавить своё объявление в избранное, ожидаем 403")
	resp := favorites.AddFavoriteExpectError(s.T(), s.loginData.Token,
		favModels.AddFavoriteRequest{AdvertisementID: ad.ID}, http.StatusForbidden)

	s.T().Logf("Ответ API: error=%q, message=%q", resp.Error, resp.Message)
	require.Equal(s.T(), "forbidden", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Cannot add own advertisement to favorites", resp.Message, "Некорректное сообщение об ошибке")

	s.T().Log("Проверены статус 403 и запрет добавления своего объявления")
}

func (s *TestSuite) TestAddFavoriteTwice() {
	s.T().Log("Авторизуем второго пользователя")
	secondUser := auth.Login(s.T(), os.Getenv("TEST_SECOND_LOGIN"), os.Getenv("TEST_SECOND_PASSWORD"))
	require.NotEqual(s.T(), s.loginData.User.Id, secondUser.User.Id, "Для теста нужны два разных аккаунта")

	s.T().Log("Первый пользователь создаёт объявление")
	adRequest := fixtures.NewAdvertisementRequest(testPhotoPath)
	ad := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, adRequest)
	defer advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)

	s.T().Log("Второй пользователь добавляет объявление в избранное первый раз")
	request := favModels.AddFavoriteRequest{AdvertisementID: ad.ID}
	favorites.AddFavorite(s.T(), secondUser.Token, request)
	defer favorites.DeleteFavorite(s.T(), secondUser.Token, ad.ID)

	s.T().Log("Добавляем то же объявление повторно, ожидаем 409")
	resp := favorites.AddFavoriteExpectError(s.T(), secondUser.Token, request, http.StatusConflict)

	s.T().Logf("Ответ API: error=%q, message=%q", resp.Error, resp.Message)
	require.Equal(s.T(), "already_exists", resp.Error, "Некорректный код ошибки")
	require.Equal(s.T(), "Advertisement is already in favorites", resp.Message, "Некорректное сообщение об ошибке")

	s.T().Log("Проверены статус 409 и запрет повторного добавления")
}
