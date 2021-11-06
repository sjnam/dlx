package dlx

import (
	"strconv"
	"strings"
)

func (mcc *MCC) InputItemNames(line string) error {
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

func (mcc *MCC) InputOptions(line string) error {
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
