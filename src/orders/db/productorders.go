package db

import "sort"

type OrderMatrix map[string]*MatrixLine

func emptyOrderMatrix(d DB) OrderMatrix {
	m := make(OrderMatrix)
	for _, p := range d.Products {
		m[p.ID] = &MatrixLine{
			Product: p,
		}
	}
	return m
}

func loadOrderMatrix(d DB, o Order) OrderMatrix {
	m := emptyOrderMatrix(d)
	for _, po := range o.Products {
		ml := m[po.ProductID]
		ml.Quantity = po.Quantity
		ml.Sliced = po.Sliced
	}

	return m
}

func (om OrderMatrix) Sum(other OrderMatrix) {
	for oid, oline := range other {
		ml := om[oid]
		ml.Quantity += oline.Quantity
		ml.Sliced += oline.Sliced
	}
}

func (om OrderMatrix) LinesByID(omitEmpty bool) []MatrixLine {
	mls := &matrixLineSorter{
		predicate: compareProductIDs,
	}
	mls.lines = make([]MatrixLine, 0, len(om))
	for _, line := range om {
		if !omitEmpty || line.Quantity > 0 {
			mls.lines = append(mls.lines, *line)
		}
	}

	sort.Sort(mls)
	return mls.lines
}

type MatrixLine struct {
	Quantity, Sliced int
	Product          Product
}

type matrixLineSorter struct {
	lines     []MatrixLine
	predicate func(i, j MatrixLine) bool
}

func (mls *matrixLineSorter) Len() int {
	return len(mls.lines)
}

func (mls *matrixLineSorter) Swap(i, j int) {
	mls.lines[i], mls.lines[j] = mls.lines[j], mls.lines[i]
}

func (mls *matrixLineSorter) Less(i, j int) bool {
	return mls.predicate(mls.lines[i], mls.lines[j])
}

func compareProductIDs(i, j MatrixLine) bool {
	return i.Product.ID < j.Product.ID
}
