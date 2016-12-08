package server

import (
	"net/http"
	"strings"
)

// deleteOrderOrProduct deletes a product or an order.
// Expected URL pattern: .../(orders|products)/<id>
func (os *OrderServer) deleteOrderOrProduct(w http.ResponseWriter, r *http.Request) {
	// Split the URL path:
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid request path", http.StatusInternalServerError)
		return
	}

	// Infer the type and an ID to be deleted:
	last := len(parts) - 1
	typ := parts[last-1]
	id := parts[last]

	// Delete:
	switch typ {
	case "orders":
		os.DB.DeleteOrderByID(id)
	case "products":
		os.DB.DeleteProductByID(id)
	}

	// Write the DB on disk:
	err := os.DB.Save()
	if err != nil {
		http.Error(w, "fail to save the db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// The current page become unavailable, so redirect to the home page:
	http.Redirect(w, r, "/"+typ, http.StatusMovedPermanently)
}
