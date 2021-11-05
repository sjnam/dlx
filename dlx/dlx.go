package dlx

import (
	"errors"
	"io"
)

const (
	infty    = int(^uint(0) >> 1)
	maxCount = infty
	root     = 0             // cl[root] is the gateway to the unsettled items
	maxLevel = 500           // at most this many options in a solution
	maxCols  = 10000         // at most this many items
	maxNodes = 100000000     // at most this many nonzero elements in the matrix
	maxLine  = 9*maxCols + 3 // a size big enough to hold all item names
)

var (
	ErrInputLineTooLong    = errors.New("input line way too long")
	ErrNoItems             = errors.New("no items")
	ErrEmptyItemName       = errors.New("empty item name")
	ErrUnknownItemName     = errors.New("unknown item name")
	ErrItemNameTooLong     = errors.New("item name too long")
	ErrDuplicateItemName   = errors.New("duplicate item name")
	ErrIllegalCharacter    = errors.New("illegal character in item name")
	ErrIllegalItemNameLine = errors.New("item name line contains | twice")
	ErrTooManyItems        = errors.New("too many items")
	ErrTooManyNodes        = errors.New("too many nodes")
	ErrTooManyLevels       = errors.New("too many levels")
	ErrPrimaryItemColored  = errors.New("primary item must be uncolored")
)

// Dancer solves exact cover problem while dancing.
type Dancer interface {
	Dance() <-chan [][]string
}

type node struct {
	up, down  int    // predecessor and successor in item list
	itm       int    // the item containing this node
	color     int    // the color specified by this node, if any
	colorName string // the color name string
}

type item struct {
	name         string // symbolic identification of the item, for printing
	prev, next   int    // neighbors of this item
	bound, slack int    // residual capacity of ths item
}

// MCC dancing links object
type MCC struct {
	nd       []node // the master list of nodes
	lastNode int    // the first node in nd that's not yet used
	cl       []item // the master list of items
	lastItm  int    // the first item in cl that's not yet used
	second   int    // boundary between primary and secondary items
}

// NewDancer generates a dancing machine.
func NewDancer(rd io.Reader) (*MCC, error) {
	d := &MCC{
		nd:     make([]node, maxNodes),
		cl:     make([]item, maxCols+2),
		second: maxCols,
	}

	if err := d.inputMatrix(rd); err != nil {
		return nil, err
	}

	return d, nil
}
