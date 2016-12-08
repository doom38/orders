package server

import (
	"net/http"
	"orders/db"
	"path"
	"strconv"
	"time"
)

// getProducts return the product list page:
func (os *OrderServer) getProducts(w http.ResponseWriter, r *http.Request) {
	os.Templater.Render(w, os.DB.Products, "page.tmpl.html", "products.tmpl.html")
}

// getProduct return the detailed product page:
func (os *OrderServer) getProduct(w http.ResponseWriter, r *http.Request) {
	// Read product ID from request path:
	d, id := path.Split(r.URL.Path)
	if d != "/products/" || id == "" {
		http.NotFound(w, r)
		return
	}

	var product db.Product
	if id != "new" {
		var found bool
		product, found = os.DB.ProductByID(id)
		if !found {
			http.NotFound(w, r)
			return
		}
	}

	m := map[string]interface{}{
		"IsNew":   id == "new",
		"Product": product,
	}
	os.Templater.Render(w, m, "page.tmpl.html", "product.tmpl.html")
	return
}

// postProduct create a new product and save the DB on the disk:
func (os *OrderServer) postProduct(w http.ResponseWriter, r *http.Request) {

	os.DB.Products = append(os.DB.Products, db.Product{
		ID:       strconv.FormatInt(time.Now().Unix(), 10),
		Name:     r.FormValue("name"),
		Slicable: r.FormValue("slicable") == "true",
	})

	err := os.DB.Save()
	if err != nil {
		http.Error(w, "fail to save the db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusMovedPermanently)
}

// patchProduct update a product and save the DB on the disk:
func (os *OrderServer) patchProduct(w http.ResponseWriter, r *http.Request) {
	d, id := path.Split(r.URL.Path)
	if d != "/products/" || id == "" {
		http.Error(w, "invalid request path", http.StatusInternalServerError)
		return
	}

	p := db.Product{
		ID:       id,
		Name:     r.FormValue("name"),
		Slicable: r.FormValue("slicable") == "true",
	}

	os.DB.UpdateProduct(p)

	err := os.DB.Save()
	if err != nil {
		http.Error(w, "fail to save the db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusMovedPermanently)
}
