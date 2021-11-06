package dlx

import (
	"strconv"
	"strings"
)

func (xcc *XCC) InputItemNames(line string) error {
	if line == "" {
		return ErrNoItems
	}

	xcc.lastItm = 1

	cl, nd := xcc.cl, xcc.nd

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if xcc.second != maxCols {
				return ErrIllegalItemNameLine
			}
			xcc.second = xcc.lastItm
			continue
		}

		if strings.ContainsAny(itm, ":|") {
			return ErrIllegalCharacter
		}

		if len(itm) > maxNameLength {
			return ErrItemNameTooLong
		}

		cl[xcc.lastItm].name = itm

		// Check for duplicate item name
		var k int
		for k = 1; cl[k].name != cl[xcc.lastItm].name; k++ {
		}
		if k < xcc.lastItm {
			return ErrDuplicateItemName
		}

		// Initialize lastItm to a new item with an empty list
		if xcc.lastItm > maxCols {
			return ErrTooManyItems
		}

		cl[xcc.lastItm-1].next = xcc.lastItm
		cl[xcc.lastItm].prev = xcc.lastItm - 1

		// nd[lastItm].itm = 0 (len)
		nd[xcc.lastItm].down = xcc.lastItm
		nd[xcc.lastItm].up = xcc.lastItm
		xcc.lastItm++
	}

	if xcc.second == maxCols {
		xcc.second = xcc.lastItm
	}

	cl[xcc.lastItm].prev = xcc.lastItm - 1
	cl[xcc.lastItm-1].next = xcc.lastItm

	cl[xcc.second].prev = xcc.lastItm
	cl[xcc.lastItm].next = xcc.second
	// this sequence works properly whether second == lastItm

	cl[root].prev = xcc.second - 1
	cl[xcc.second-1].next = root

	xcc.lastNode = xcc.lastItm // reserve all the header nodes and the first spacer
	// we have nd[lastNode].itm=0 in the first spacer

	return nil
}

func (xcc *XCC) InputOptions(line string) error {
	var (
		cl      = xcc.cl
		nd      = xcc.nd
		options = 0            // options seen so far
		i       = xcc.lastNode // remember the spacer at the left of this option
		pp      = false
	)

	for _, opt := range strings.Fields(line) {
		if len(opt) > maxNameLength {
			return ErrItemNameTooLong
		}
		if opt[0] == ':' {
			return ErrEmptyItemName
		}
		name := strings.Split(opt, ":")
		cl[xcc.lastItm].name = name[0]

		// Create a node for the item named in opt
		var k int
		for k = 0; cl[k].name != cl[xcc.lastItm].name; k++ {
		}
		if k == xcc.lastItm {
			return ErrUnknownItemName
		}
		if nd[k].color >= i { // aux field
			return ErrDuplicateItemName
		}
		xcc.lastNode++
		if xcc.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		nd[xcc.lastNode].itm = k
		if k < xcc.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		nd[k].itm++                // len field; store the new length of the list
		nd[k].color = xcc.lastNode // aux field

		r := nd[k].up // the "bottom" node of the item list
		nd[k].up = xcc.lastNode
		nd[r].down = xcc.lastNode
		nd[xcc.lastNode].up = r
		nd[xcc.lastNode].down = k

		if len(name) == 1 {
			nd[xcc.lastNode].color = 0
			nd[xcc.lastNode].colorName = ""
		} else if k >= xcc.second {
			c, _ := strconv.ParseInt(name[1], 36, 0)
			nd[xcc.lastNode].color = int(c)
			nd[xcc.lastNode].colorName = ":" + name[1]
		} else {
			return ErrPrimaryItemColored
		}
	}

	if !pp {
		for xcc.lastNode > i {
			// Remove lastNode from its item list
			k := nd[xcc.lastNode].itm
			nd[k].itm--
			nd[k].color = i - 1
			q := nd[xcc.lastNode].up
			r := nd[xcc.lastNode].down
			nd[q].down = r
			nd[r].up = q
			xcc.lastNode--
		}
	} else {
		nd[i].down = xcc.lastNode
		xcc.lastNode++ // create the next spacer
		if xcc.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		options++
		nd[xcc.lastNode].up = i + 1
		nd[xcc.lastNode].itm = -options
	}

	return nil
}
