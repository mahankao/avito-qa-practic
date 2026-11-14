package models

type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	Photos      []string
}

// ErrorResponse описывает ответ сервера с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details"`
}
