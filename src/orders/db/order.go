package db

// Order is a list quantity-product pair and some other metadata
type Order struct {
	ID           string         `json:"id"`
	CustomerName string         `json:"customer-name"`
	Comment      string         `json:"comment"`
	DueHour      string         `json:"due-hour"`
	WithBag      bool           `json:"with-bag"`
	Products     []ProductOrder `json:"products"`
}
