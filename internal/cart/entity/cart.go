package entity

type Cart struct {
	Products []Products `json:"products"`
}

type Products struct {
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}
