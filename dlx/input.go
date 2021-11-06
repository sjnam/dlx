package dlx

import (
	"strconv"
	"strings"
)

// XC

func (xc *XC) inputItemNames(line string) error {
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

func (xc *XC) inputOptions(line string) error {
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

// XCC

func (xcc *XCC) inputItemNames(line string) error {
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

func (xcc *XCC) inputOptions(line string) error {
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

// MCC

func (mcc *MCC) inputItemNames(line string) error {
	if line == "" {
		return ErrNoItems
	}

	mcc.lastItm = 1

	cl, nd := mcc.cl, mcc.nd

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if mcc.second != maxCols {
				return ErrIllegalItemNameLine
			}
			mcc.second = mcc.lastItm
			continue
		}

		if len(itm) > maxNameLength {
			return ErrItemNameTooLong
		}

		q, r := 1, 1
		bn := strings.Split(itm, "|")
		if len(bn) == 1 {
			cl[mcc.lastItm].name = itm
		} else {
			cl[mcc.lastItm].name = bn[1]
			bounds := strings.Split(bn[0], ":")
			if len(bounds) == 1 {
				q, _ = strconv.Atoi(bounds[0])
				r = q
			} else {
				r, _ = strconv.Atoi(bounds[0])
				q, _ = strconv.Atoi(bounds[1])
			}
		}
		cl[mcc.lastItm].bound = q
		cl[mcc.lastItm].slack = q - r

		// Check for duplicate item name
		var k int
		for k = 1; cl[k].name != cl[mcc.lastItm].name; k++ {
		}
		if k < mcc.lastItm {
			return ErrDuplicateItemName
		}

		// Initialize lastItm to a new item with an empty list
		if mcc.lastItm > maxCols {
			return ErrTooManyItems
		}

		cl[mcc.lastItm-1].next = mcc.lastItm
		cl[mcc.lastItm].prev = mcc.lastItm - 1

		// nd[lastItm].itm = 0 (len)
		nd[mcc.lastItm].down = mcc.lastItm
		nd[mcc.lastItm].up = mcc.lastItm
		mcc.lastItm++
	}

	if mcc.second == maxCols {
		mcc.second = mcc.lastItm
	}

	cl[mcc.lastItm].prev = mcc.lastItm - 1
	cl[mcc.lastItm-1].next = mcc.lastItm

	cl[mcc.second].prev = mcc.lastItm
	cl[mcc.lastItm].next = mcc.second
	// this sequence works properly whether second == lastItm

	cl[root].prev = mcc.second - 1
	cl[mcc.second-1].next = root

	mcc.lastNode = mcc.lastItm // reserve all the header nodes and the first spacer
	// we have nd[lastNode].itm=0 in the first spacer

	return nil
}

func (mcc *MCC) inputOptions(line string) error {
	var (
		cl      = mcc.cl
		nd      = mcc.nd
		options = 0            // options seen so far
		i       = mcc.lastNode // remember the spacer at the left of this option
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
		cl[mcc.lastItm].name = name[0]

		// Create a node for the item named in opt
		var k int
		for k = 0; cl[k].name != cl[mcc.lastItm].name; k++ {
		}
		if k == mcc.lastItm {
			return ErrUnknownItemName
		}
		if nd[k].color >= i { // aux field
			return ErrDuplicateItemName
		}
		mcc.lastNode++
		if mcc.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		nd[mcc.lastNode].itm = k
		if k < mcc.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		nd[k].itm++                // len field; store the new length of the list
		nd[k].color = mcc.lastNode // aux field

		r := nd[k].up // the "bottom" node of the item list
		nd[k].up = mcc.lastNode
		nd[r].down = mcc.lastNode
		nd[mcc.lastNode].up = r
		nd[mcc.lastNode].down = k

		if len(name) == 1 {
			nd[mcc.lastNode].color = 0
			nd[mcc.lastNode].colorName = ""
		} else if k >= mcc.second {
			c, _ := strconv.ParseInt(name[1], 36, 0)
			nd[mcc.lastNode].color = int(c)
			nd[mcc.lastNode].colorName = ":" + name[1]
		} else {
			return ErrPrimaryItemColored
		}
	}

	if !pp {
		for mcc.lastNode > i {
			// Remove lastNode from its item list
			k := nd[mcc.lastNode].itm
			nd[k].itm--
			nd[k].color = i - 1
			q := nd[mcc.lastNode].up
			r := nd[mcc.lastNode].down
			nd[q].down = r
			nd[r].up = q
			mcc.lastNode--
		}
	} else {
		nd[i].down = mcc.lastNode
		mcc.lastNode++ // create the next spacer
		if mcc.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		options++
		nd[mcc.lastNode].up = i + 1
		nd[mcc.lastNode].itm = -options
	}

	return nil
}
