package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// Load a DB from a JSON file.
func Load(path string) (db DB, err error) {
	db.filePath = path

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &db)

	return
}

// DB is the storage main object.
type DB struct {
	Products []Product `json:"products"`
	Orders   []Order   `json:"orders"`

	filePath string
}

// Save writes the DB on the disk.
func (d *DB) Save() (err error) {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(d.filePath, b, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("Can't save the DB at path %q: %v", d.filePath, err)
		return
	}
	return
}

// ProductByID looks for a product in the DB using its ID.
// 'found' is false when no product with this id is found in the DB.
func (d *DB) ProductByID(id string) (p Product, found bool) {
	for _, product := range d.Products {
		if product.ID == id {
			return product, true
		}
	}
	return
}

// OrderByID looks for an order in the DB using its ID.
// 'found' is false when no order with this id is found in the DB.
func (d *DB) OrderByID(id string) (o Order, found bool) {
	for _, order := range d.Orders {
		if order.ID == id {
			return order, true
		}
	}
	return
}

// OrdersByDueHour returns the order list sorted using the due hour.
func (d *DB) OrdersByDueHour() []Order {
	o := make([]Order, len(d.Orders))
	for i, j := range d.Orders {
		o[i] = j
	}

	sorter := &orderSorter{
		orders:    o,
		predicate: compareOrderDueHour,
	}
	sort.Sort(sorter)

	return sorter.orders
}

// OrderMatrix returns the order content in a way easy to render using a
// template.
func (d *DB) OrderMatrix(orderID string) OrderMatrix {
	o, _ := d.OrderByID(orderID)
	return loadOrderMatrix(*d, o)
}

// TotalOrder sums all the DB orders and returns a matrix
func (d *DB) TotalOrder() OrderMatrix {
	to := emptyOrderMatrix(*d)
	for _, order := range d.Orders {
		m := d.OrderMatrix(order.ID)
		to.Sum(m)
	}
	return to
}

// Update a product in the DB
func (d *DB) UpdateProduct(p Product) bool {
	for i, op := range d.Products {
		if op.ID == p.ID {
			d.Products[i] = p
			return true
		}
	}
	return false
}

// Update an order in the DB
func (d *DB) UpdateOrder(o Order) bool {
	for i, oo := range d.Orders {
		if oo.ID == o.ID {
			d.Orders[i] = o
			return true
		}
	}
	return false
}

// DeleteOrderByID remove the order denoted by this ID
func (d *DB) DeleteOrderByID(id string) (deleted bool) {
	list := make([]Order, 0, len(d.Orders))
	for _, o := range d.Orders {
		if o.ID != id {
			list = append(list, o)
		} else {
			deleted = true
		}
	}
	d.Orders = list
	return
}

// DeleteProductByID remove the product denoted by this ID
func (d *DB) DeleteProductByID(id string) (deleted bool) {
	list := make([]Product, 0, len(d.Products))
	for _, p := range d.Products {
		if p.ID != id {
			list = append(list, p)
		} else {
			deleted = true
		}
	}
	d.Products = list
	return
}
