package db

// orderSorter sorts an order list according to a predicate.
type orderSorter struct {
	orders    []Order
	predicate func(i, j Order) bool
}

func (os *orderSorter) Len() int {
	return len(os.orders)
}

func (os *orderSorter) Swap(i, j int) {
	os.orders[i], os.orders[j] = os.orders[j], os.orders[i]
}

func (os *orderSorter) Less(i, j int) bool {
	return os.predicate(os.orders[i], os.orders[j])
}

func compareOrderDueHour(i, j Order) bool {
	return i.DueHour < j.DueHour
}
