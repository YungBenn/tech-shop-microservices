package entity

type ProductData struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Price       string   `json:"price"`
	Tag         []string `json:"tag"`
	Discount    string   `json:"discount"`
	Image       []string `json:"image"`
	Description string   `json:"description"`
	CreatedBy   string   `json:"created_by"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}
