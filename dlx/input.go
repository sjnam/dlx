package dlx

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

const maxNameLength = 32

func (m *MCC) inputMatrix(rd io.Reader) error {
	var line string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > maxLine {
			return ErrInputLineTooLong
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' {
			// bypass comment or blank line
			continue
		}
		m.lastItm = 1
		break
	}

	if m.lastItm == 0 {
		return ErrNoItems
	}

	if err := m.inputItemNames(line); err != nil {
		return err
	}

	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > maxLine {
			return ErrInputLineTooLong
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' {
			// bypass comment or blank line
			continue
		}

		if err := m.inputOptions(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (m *MCC) inputItemNames(line string) error {
	cl, nd := m.cl, m.nd

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if m.second != maxCols {
				return ErrIllegalItemNameLine
			}
			m.second = m.lastItm
			continue
		}

		if len(itm) > maxNameLength {
			return ErrItemNameTooLong
		}

		q, r := 1, 1
		bn := strings.Split(itm, "|")
		if len(bn) == 1 {
			cl[m.lastItm].name = itm
		} else {
			cl[m.lastItm].name = bn[1]
			bounds := strings.Split(bn[0], ":")
			if len(bounds) == 1 {
				q, _ = strconv.Atoi(bounds[0])
				r = q
			} else {
				r, _ = strconv.Atoi(bounds[0])
				q, _ = strconv.Atoi(bounds[1])
			}
		}
		cl[m.lastItm].bound = q
		cl[m.lastItm].slack = q - r

		// Check for duplicate item name
		var k int
		for k = 1; cl[k].name != cl[m.lastItm].name; k++ {
		}
		if k < m.lastItm {
			return ErrDuplicateItemName
		}

		// Initialize lastItm to a new item with an empty list
		if m.lastItm > maxCols {
			return ErrTooManyItems
		}

		cl[m.lastItm-1].next = m.lastItm
		cl[m.lastItm].prev = m.lastItm - 1

		// nd[lastItm].itm = 0 (len)
		nd[m.lastItm].down = m.lastItm
		nd[m.lastItm].up = m.lastItm
		m.lastItm++
	}

	if m.second == maxCols {
		m.second = m.lastItm
	}

	cl[m.lastItm].prev = m.lastItm - 1
	cl[m.lastItm-1].next = m.lastItm

	cl[m.second].prev = m.lastItm
	cl[m.lastItm].next = m.second
	// this sequence works properly whether second == lastItm

	cl[root].prev = m.second - 1
	cl[m.second-1].next = root

	m.lastNode = m.lastItm // reserve all the header nodes and the first spacer
	// we have nd[lastNode].itm=0 in the first spacer

	return nil
}

func (m *MCC) inputOptions(line string) error {
	var (
		cl      = m.cl
		nd      = m.nd
		options = 0          // options seen so far
		i       = m.lastNode // remember the spacer at the left of this option
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
		cl[m.lastItm].name = name[0]

		// Create a node for the item named in opt
		var k int
		for k = 0; cl[k].name != cl[m.lastItm].name; k++ {
		}
		if k == m.lastItm {
			return ErrUnknownItemName
		}
		if nd[k].color >= i { // aux field
			return ErrDuplicateItemName
		}
		m.lastNode++
		if m.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		nd[m.lastNode].itm = k
		if k < m.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		nd[k].itm++              // len field; store the new length of the list
		nd[k].color = m.lastNode // aux field

		r := nd[k].up // the "bottom" node of the item list
		nd[k].up = m.lastNode
		nd[r].down = m.lastNode
		nd[m.lastNode].up = r
		nd[m.lastNode].down = k

		if len(name) == 1 {
			nd[m.lastNode].color = 0
			nd[m.lastNode].colorName = ""
		} else if k >= m.second {
			c, _ := strconv.ParseInt(name[1], 36, 0)
			nd[m.lastNode].color = int(c)
			nd[m.lastNode].colorName = ":" + name[1]
		} else {
			return ErrPrimaryItemColored
		}
	}

	if !pp {
		for m.lastNode > i {
			// Remove lastNode from its item list
			k := nd[m.lastNode].itm
			nd[k].itm--
			nd[k].color = i - 1
			q := nd[m.lastNode].up
			r := nd[m.lastNode].down
			nd[q].down = r
			nd[r].up = q
			m.lastNode--
		}
	} else {
		nd[i].down = m.lastNode
		m.lastNode++ // create the next spacer
		if m.lastNode == maxNodes {
			return ErrTooManyNodes
		}
		options++
		nd[m.lastNode].up = i + 1
		nd[m.lastNode].itm = -options
	}

	return nil
}
