# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
API_URL=https://testboard.avito.com/api/v1
DEBUG=false

TEST_LOGIN=first_user_email
TEST_PASSWORD=first_user_password

TEST_SECOND_LOGIN=second_user_email
TEST_SECOND_PASSWORD=second_user_password
```
* Для запуска нужны два разных аккаунта Testboard.
Первый создаёт объявления, второй добавляет их в избранное.
Оба аккаунта необходимо зарегистрировать заранее

## Подготовка тестовых данных

Для существующих тестов поиска и собственных объявлений
у первого пользователя должно быть объявление с названием
«Цветы», непустым описанием и хотя бы одной фотографией.

Новые тесты создания и избранного сами создают необходимые
объявления и удаляют их после выполнения.

Тестовая фотография хранится в:
tests/scenarios/advertisementCreate/testdata/photo.jpg

Пользователи должны иметь свободное место для создания
объявлений и добавления записей в избранное.



8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -v ./...
```
## Проверка лабораторной работы

Тесты создания:
go test -v ./tests/scenarios/advertisementCreate/

Тесты избранного:
go test -v ./tests/scenarios/favorites/

Все тесты три раза подряд:
go test ./... -count=3

Линтер — использована версия 2.6.1:
golangci-lint run

9. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -v ./tests/scenarios/myAdvertisement/
```

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```