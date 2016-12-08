package server

import (
	"assets"
	"fmt"
	"log"
	"mime"
	"net/http"
	"orders/db"
	"orders/tmpl"
	"path"
	"strings"
)

// OrderServer serves all the applications requests.
type OrderServer struct {
	DB        *db.DB         // Storage
	Templater tmpl.Templater // Template engine
}

func (os *OrderServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rule := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	log.Println(rule)

	switch {
	// assets requests:
	case strings.HasPrefix(rule, "GET /assets/"):
		handler := http.HandlerFunc(serveAssets)
		http.StripPrefix("/assets/", handler).ServeHTTP(w, r)

	// Home, currently no home page, just a redirect:
	case rule == "GET /":
		http.Redirect(w, r, "/orders", http.StatusMovedPermanently)

	// Products requests:
	case rule == "GET /products":
		os.getProducts(w, r)
	case strings.HasPrefix(rule, "GET /products/"):
		os.getProduct(w, r)
	case rule == "POST /products":
		os.postProduct(w, r)
	case strings.HasPrefix(rule, "POST /products/"):
		os.patchProduct(w, r)

	// Orders requests:
	case rule == "GET /orders":
		os.getOrders(w, r)
	case strings.HasPrefix(rule, "GET /orders/"):
		os.getOrder(w, r)
	case rule == "POST /orders":
		os.postOrder(w, r)
	case strings.HasPrefix(rule, "POST /orders/"):
		os.patchOrder(w, r)

	// Delete operations, not REST compliant but easiest in my case.
	case strings.HasPrefix(rule, "POST /delete/"):
		os.deleteOrderOrProduct(w, r)

	// Otherwise return 404 error:
	default:
		log.Println("404 ->", rule)
		http.NotFound(w, r)
		return
	}
}

// serveAssets serves assets embedded in the package 'assets'.
// 'asset' is generated using:
// https://github.com/jteeuwen/go-bindata
func serveAssets(w http.ResponseWriter, r *http.Request) {
	data, err := assets.Asset(r.URL.Path)
	if err != nil {
		log.Println("404 -> ASSET", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	// Infer the asset type from its file extension, to set the 'Content-Type'
	// header:
	ext := path.Ext(r.URL.Path)
	w.Header()["Content-Type"] = []string{mime.TypeByExtension(ext)}
	w.Write(data)
}
