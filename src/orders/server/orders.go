package server

import (
	"net/http"
	"orders/db"
	"path"
	"strconv"
	"time"
)

// getOrders return the order list page:
func (os *OrderServer) getOrders(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["Sums"] = os.DB.TotalOrder()

	orders := make([]interface{}, len(os.DB.Orders))
	m["Orders"] = orders
	for i, o := range os.DB.OrdersByDueHour() {
		orders[i] = map[string]interface{}{
			"ID":           o.ID,
			"CustomerName": o.CustomerName,
			"DueHour":      o.DueHour,
			"WithBag":      o.WithBag,
			"Comment":      o.Comment,
			"Matrix":       os.DB.OrderMatrix(o.ID).LinesByID(true),
		}
	}

	os.Templater.Render(w, m, "page.tmpl.html", "orders.tmpl.html")
}

// getOrder return the detailed order page:
func (os *OrderServer) getOrder(w http.ResponseWriter, r *http.Request) {
	d, id := path.Split(r.URL.Path)
	if d != "/orders/" || id == "" {
		http.NotFound(w, r)
		return
	}

	var order db.Order
	order.ID = id
	if id != "new" {
		var found bool
		order, found = os.DB.OrderByID(id)
		if !found {
			http.NotFound(w, r)
			return
		}
	}

	m := map[string]interface{}{
		"Order":  order,
		"IsNew":  order.ID == "new",
		"Matrix": os.DB.OrderMatrix(order.ID).LinesByID(false),
	}
	os.Templater.Render(w, m, "page.tmpl.html", "order.tmpl.html")
}

// postOrder create a new order and save the DB on the disk:
func (os *OrderServer) postOrder(w http.ResponseWriter, r *http.Request) {
	// Read the form values:
	o := parseOrder(r, *os.DB)

	// Generate a uniq ID (UNIX timestamp)
	// /!\ Unsecure but easy to implement
	o.ID = strconv.FormatInt(time.Now().Unix(), 10)

	// Update the order list:
	os.DB.Orders = append(os.DB.Orders, o)

	// Save the DB
	err := os.DB.Save()
	if err != nil {
		http.Error(w, "fail to save the db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the updated order list:
	http.Redirect(w, r, "/orders", http.StatusMovedPermanently)
}

// patchOrder update an order and save the DB on the disk:
func (os *OrderServer) patchOrder(w http.ResponseWriter, r *http.Request) {
	// Read order ID from request path
	d, id := path.Split(r.URL.Path)
	if d != "/orders/" || id == "" {
		http.Error(w, "invalid request path", http.StatusInternalServerError)
		return
	}

	// Read the form values:
	newOrder := parseOrder(r, *os.DB)
	newOrder.ID = id

	// Update the order in the DB:
	os.DB.UpdateOrder(newOrder)

	// Save the DB:
	err := os.DB.Save()
	if err != nil {
		http.Error(w, "fail to save the db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the updated order list:
	http.Redirect(w, r, "/orders", http.StatusMovedPermanently)
}

// parseOrder reads/parse the form values from the provided request:
func parseOrder(r *http.Request, d db.DB) db.Order {
	var o db.Order
	o.CustomerName = r.FormValue("customer-name")
	o.DueHour = r.FormValue("due-hour")
	o.WithBag = r.FormValue("with-bag") == "true"
	o.Comment = r.FormValue("comment")

	for _, p := range d.Products {
		q := asInt(r.FormValue(p.ID + ".quantity"))
		if q > 0 {
			o.Products = append(o.Products, db.ProductOrder{
				ProductID: p.ID,
				Quantity:  q,
				Sliced:    asInt(r.FormValue(p.ID + ".sliced")),
			})
		}
	}

	return o
}

// asInt is a convertion method that returns 0 when 'val' is empty or when 'val'
// isn't an integer.
func asInt(val string) int {
	i64, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0
	}
	return int(i64)
}
