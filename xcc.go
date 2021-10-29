package dlx

import (
	"math/rand"
	"time"
)

const (
	root     = 0             // cl[root] is the gateway to the unsettled items
	maxCols  = 100000        // at most this many items
	maxNodes = 10000000      // at most this many nonzero elements in the matrix
	maxLine  = 9*maxCols + 3 // a size big enough to hold all item names
	maxLevel = 5000          // at most this many options in a solution
)

type node struct {
	up, down int    // predecessor and successor in item list
	itm      int    // the item containing this node
	color    int    // the color specified by this node, if any
	scolor   string // color name
}

type item struct {
	name       string // symbolic identification of the item, for printing
	prev, next int    // neighbors of this item
}

type XCC struct {
	nd       []node // the master list of nodes
	lastNode int    // the first node in nd that's not yet used
	cl       []item // the master list of items
	lastItm  int    // the first item in cl that's not yet used
	second   int    // boundary between primary and secondary items
	options  int    // options seen so far
	choice   []int  // the node chosen on each level
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewXCC() *XCC {
	return &XCC{
		nd:     make([]node, maxNodes),
		cl:     make([]item, maxCols+2),
		second: maxCols,
		choice: make([]int, maxLevel),
	}
}
