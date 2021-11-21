package dlx

import (
	"context"
	"io"
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

	maxNameLength = 128 // max item name length
	chunkSize     = 256
)

type Option []string

type Result struct {
	Solutions <-chan []Option
	Heartbeat <-chan string
}

type Dancer interface {
	Dance(reader io.Reader) Result
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
	ctx           context.Context
	nd            []node // the master list of nodes
	lastNode      int    // the first node in nd that's not yet used
	cl            []item // the master list of items
	lastItm       int    // the first item in cl that's not yet used
	second        int    // boundary between primary and secondary items
	options       int    // options seen so far
	updates       uint64 // update count
	cleansings    uint64 // cleansing count
	Debug         bool   // info/debug message
	PulseInterval time.Duration
}

// NewMCC Wake me up before you Go Go
func NewMCC() *MCC {
	return &MCC{
		nd:            make([]node, chunkSize),
		cl:            make([]item, chunkSize),
		second:        maxCols,
		ctx:           context.Background(),
		PulseInterval: time.Hour,
	}
}

func (m *MCC) WithContext(ctx context.Context) *MCC {
	if ctx == nil {
		panic("nil context")
	}
	d2 := new(MCC)
	*d2 = *m
	d2.ctx = ctx
	return d2
}
