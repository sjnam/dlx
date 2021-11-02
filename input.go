package dlx

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const MaxNameLength = 32

func (d *DLX) inputMatrix(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > maxLine {
			return fmt.Errorf("input line way too long")
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' {
			// bypass comment or blank line
			continue
		}
		d.lastItm = 1
		break
	}

	if d.lastItm == 0 {
		return fmt.Errorf("no items")
	}

	if err := d.inputItemNames(line); err != nil {
		return err
	}

	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > maxLine {
			return fmt.Errorf("option line too long")
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' {
			// bypass comment or blank line
			continue
		}

		if err := d.inputOptions(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (d *DLX) inputItemNames(line string) error {
	cl, nd := d.cl, d.nd

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if d.second != maxCols {
				return fmt.Errorf("item name line contains | twice")
			}
			d.second = d.lastItm
			continue
		}

		if strings.ContainsAny(itm, ":|") {
			return fmt.Errorf("illegal character in item name")
		}

		if len(itm) > MaxNameLength {
			return fmt.Errorf("item name too long")
		}

		cl[d.lastItm].name = itm

		// Check for duplicate item name
		var k int
		for k = 1; cl[k].name != cl[d.lastItm].name; k++ {
		}
		if k < d.lastItm {
			return fmt.Errorf("duplicate item name")
		}

		// Initialize lastItm to a new item with an empty list
		if d.lastItm > maxCols {
			return fmt.Errorf("too many items")
		}

		cl[d.lastItm-1].next = d.lastItm
		cl[d.lastItm].prev = d.lastItm - 1

		// nd[lastItm].itm = 0 (len)
		nd[d.lastItm].down = d.lastItm
		nd[d.lastItm].up = d.lastItm
		d.lastItm++
	}

	if d.second == maxCols {
		d.second = d.lastItm
	}

	cl[d.lastItm].prev = d.lastItm - 1
	cl[d.lastItm-1].next = d.lastItm

	cl[d.second].prev = d.lastItm
	cl[d.lastItm].next = d.second
	// this sequence works properly whether second == lastItm

	cl[root].prev = d.second - 1
	cl[d.second-1].next = root

	d.lastNode = d.lastItm // reserve all the header nodes and the first spacer
	// we have nd[lastNode].itm=0 in the first spacer

	return nil
}

func (d *DLX) inputOptions(line string) error {
	var (
		cl      = d.cl
		nd      = d.nd
		options = 0          // options seen so far
		i       = d.lastNode // remember the spacer at the left of this option
		pp      = false
	)

	for _, opt := range strings.Fields(line) {
		if len(opt) > MaxNameLength {
			return fmt.Errorf("item name too long")
		}
		if opt[0] == ':' {
			return fmt.Errorf("empty item name")
		}
		name := strings.Split(opt, ":")
		cl[d.lastItm].name = name[0]

		// Create a node for the item named in opt
		var k int
		for k = 0; cl[k].name != cl[d.lastItm].name; k++ {
		}
		if k == d.lastItm {
			return fmt.Errorf("unknown item name")
		}
		if nd[k].color >= i { // aux field
			return fmt.Errorf("duplicate item name in this option")
		}
		d.lastNode++
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		nd[d.lastNode].itm = k
		if k < d.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		nd[k].itm++              // len field; store the new length of the list
		nd[k].color = d.lastNode // aux field

		r := nd[k].up // the "bottom" node of the item list
		nd[k].up = d.lastNode
		nd[r].down = d.lastNode
		nd[d.lastNode].up = r
		nd[d.lastNode].down = k

		if len(name) == 1 {
			nd[d.lastNode].color = 0
			nd[d.lastNode].colorName = ""
		} else if k >= d.second {
			c, _ := strconv.ParseInt(name[1], 36, 0)
			nd[d.lastNode].color = int(c)
			nd[d.lastNode].colorName = ":" + name[1]
		} else {
			return fmt.Errorf("primary item must be uncolored")
		}
	}

	if !pp {
		for d.lastNode > i {
			// Remove lastNode from its item list
			k := nd[d.lastNode].itm
			nd[k].itm--
			nd[k].color = i - 1
			q := nd[d.lastNode].up
			r := nd[d.lastNode].down
			nd[q].down = r
			nd[r].up = q
			d.lastNode--
		}
	} else {
		nd[i].down = d.lastNode
		d.lastNode++ // create the next spacer
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		options++
		nd[d.lastNode].up = i + 1
		nd[d.lastNode].itm = -options
	}

	return nil
}

func (d *DLX) InitialContent() {
	cl, nd := d.cl, d.nd
	for i := 0; i < d.lastItm; i++ {
		fmt.Printf("%2d[%2s,%2d,%2d] ",
			i, cl[i].name, cl[i].prev, cl[i].next)
	}
	fmt.Println()
	for i := 0; i <= d.lastNode; i++ {
		if i%d.lastItm == 0 {
			fmt.Println()
		}
		fmt.Printf("%2d[%2d,%2d,%2d] ",
			i, nd[i].itm, nd[i].up, nd[i].down)
	}
	fmt.Println()
}
