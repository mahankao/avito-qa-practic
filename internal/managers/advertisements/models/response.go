package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Photo struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	SortOrder int    `json:"sort_order"`
	URL       string `json:"url"`
}

type Advertisement struct {
	Title       string  `json:"title"`
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Photos      []Photo `json:"photos"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	UserID      string  `json:"user_id"`
}
