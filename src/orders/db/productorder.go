package db

// ProductOrder is a product ID & quantity pair
type ProductOrder struct {
	ProductID string `json:"product-id"`
	Quantity  int    `json:"quantity"`
	Sliced    int    `json:"sliced"`
}
