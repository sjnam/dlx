package dlx

import (
	"math/rand"
	"time"
)

const (
	infty    = int(^uint(0) >> 1)
	maxCount = infty
	root     = 0             // cl[root] is the gateway to the unsettled items
	maxLevel = 500           // at most this many options in a solution
	maxCols  = 10000         // at most this many items
	maxNodes = 100000000     // at most this many nonzero elements in the matrix
	maxLine  = 9*maxCols + 3 // a size big enough to hold all item names

	maxNameLength = 32 // max item name length
)

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

// Dancer dancing links object
type Dancer struct {
	nd       []node // the master list of nodes
	lastNode int    // the first node in nd that's not yet used
	cl       []item // the master list of items
	lastItm  int    // the first item in cl that's not yet used
	second   int    // boundary between primary and secondary items
	options  int    // options seen so far
	Info     bool   // info/debug message
}

// NewDancer Wake me up before you Go Go
func NewDancer() *Dancer {
	rand.Seed(time.Now().UnixNano())
	return &Dancer{
		nd:     make([]node, 1), // maxNodes
		cl:     make([]item, 1), // maxCols+2
		second: maxCols,
	}
}
