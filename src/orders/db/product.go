package db

// Product is identifiable and has a name
type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slicable bool   `json:"slicable"`
}
