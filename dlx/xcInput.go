package dlx

import "strings"

func (xc *XC) InputItemNames(line string) error {
	if line == "" {
		return ErrNoItems
	}

	xc.lastItm = 1

	cl, nd := xc.cl, xc.nd

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if xc.second != maxCols {
				return ErrIllegalItemNameLine
			}
			xc.second = xc.lastItm
			continue
		}

		if strings.ContainsAny(itm, ":|") {
			return ErrIllegalCharacter
		}

		if len(itm) > maxNameLength {
			return ErrItemNameTooLong
		}

		cl[xc.lastItm].name = itm

		// Check for duplicate item name
		var k int
		for k = 1; cl[k].name != cl[xc.lastItm].name; k++ {
		}
		if k < xc.lastItm {
			return ErrDuplicateItemName
		}

		// Initialize lastItm to a new item with an empty list
		if xc.lastItm > maxCols {
			return ErrTooManyItems
		}

		cl[xc.lastItm-1].next = xc.lastItm
		cl[xc.lastItm].prev = xc.lastItm - 1

		// nd[lastItm].itm = 0 (len)
		nd[xc.lastItm].down = xc.lastItm
		nd[xc.lastItm].up = xc.lastItm
		xc.lastItm++
	}

	if xc.second == maxCols {
		xc.second = xc.lastItm
	}

	cl[xc.lastItm].prev = xc.lastItm - 1
	cl[xc.lastItm-1].next = xc.lastItm

	cl[xc.second].prev = xc.lastItm
	cl[xc.lastItm].next = xc.second
	// this sequence works properly whether second == lastItm

	cl[root].prev = xc.second - 1
	cl[xc.second-1].next = root

	xc.lastNode = xc.lastItm // reserve all the header nodes and the first spacer
	// we have nd[lastNode].itm=0 in the first spacer

	return nil
}

func (xc *XC) InputOptions(line string) error {
	var (
		cl      = xc.cl
		nd      = xc.nd
		options = 0           // options seen so far
		i       = xc.lastNode // remember the spacer at the left of this option
		pp      = false
	)

	for _, opt := range strings.Fields(line) {
		if len(opt) > maxNameLength {
			return ErrItemNameTooLong
		}

		cl[xc.lastItm].name = opt

		// Create a node for the item named in opt
		var k int
		for k = 0; cl[k].name != cl[xc.lastItm].name; k++ {
		}
		if k == xc.lastItm {
			return ErrUnknownItemName
		}
		if nd[k].color >= i { // aux field
			return ErrDuplicateItemName
		}
		xc.lastNode++
		if xc.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		nd[xc.lastNode].itm = k
		if k < xc.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		nd[k].itm++               // len field; store the new length of the list
		nd[k].color = xc.lastNode // aux field

		r := nd[k].up // the "bottom" node of the item list
		nd[k].up = xc.lastNode
		nd[r].down = xc.lastNode
		nd[xc.lastNode].up = r
		nd[xc.lastNode].down = k
	}

	if !pp {
		for xc.lastNode > i {
			// Remove lastNode from its item list
			k := nd[xc.lastNode].itm
			nd[k].itm--
			nd[k].color = i - 1
			q := nd[xc.lastNode].up
			r := nd[xc.lastNode].down
			nd[q].down = r
			nd[r].up = q
			xc.lastNode--
		}
	} else {
		nd[i].down = xc.lastNode
		xc.lastNode++ // create the next spacer
		if xc.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		options++
		nd[xc.lastNode].up = i + 1
		nd[xc.lastNode].itm = -options
	}

	return nil
}
