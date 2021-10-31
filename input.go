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
	if maxNodes <= 2*maxCols {
		return fmt.Errorf("recompile me: maxNodes must exceed twice maxCols")
	} // every item will want a header node and at least one other node

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

		d.cl[d.lastItm].name = itm

		// Check for duplicate item name
		var k int
		for k = 1; d.cl[k].name != d.cl[d.lastItm].name; k++ {
		}
		if k < d.lastItm {
			return fmt.Errorf("duplicate item name")
		}

		// Initialize lastItm to a new item with an empty list
		if d.lastItm > maxCols {
			return fmt.Errorf("too many items")
		}

		d.cl[d.lastItm-1].next = d.lastItm
		d.cl[d.lastItm].prev = d.lastItm - 1

		// d.nd[lastItm].itm = 0 (len)
		d.nd[d.lastItm].down = d.lastItm
		d.nd[d.lastItm].up = d.lastItm
		d.lastItm++
	}

	if d.second == maxCols {
		d.second = d.lastItm
	}

	d.cl[d.lastItm].prev = d.lastItm - 1
	d.cl[d.lastItm-1].next = d.lastItm

	d.cl[d.second].prev = d.lastItm
	d.cl[d.lastItm].next = d.second
	// this sequence works properly whether second == lastItm

	d.cl[root].prev = d.second - 1
	d.cl[d.second-1].next = root

	d.lastNode = d.lastItm // reserve all the header nodes and the first spacer
	// we have nd[lastNode].itm=0 in the first spacer

	return nil
}

func (d *DLX) inputOptions(line string) error {
	var (
		pp = false
		i  = d.lastNode // remember the spacer at the left of this option
	)

	for _, opt := range strings.Fields(line) {
		if len(opt) > MaxNameLength {
			return fmt.Errorf("item name too long")
		}

		name := strings.Split(opt, ":")

		// Create a node for the item named in opt
		var k int
		for k = 0; d.cl[k].name != name[0]; k++ {
		}
		if k == d.lastItm {
			return fmt.Errorf("unknown item name")
		}
		if d.nd[k].color >= i {
			return fmt.Errorf("duplicate item name in this option")
		}
		d.lastNode++
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		d.nd[d.lastNode].itm = k
		if k < d.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		t := d.nd[k].itm + 1
		d.nd[k].itm = t            // store the new length of the list
		d.nd[k].color = d.lastNode // aux field

		r := d.nd[k].up // the "bottom" node of the item list
		d.nd[k].up = d.lastNode
		d.nd[r].down = d.lastNode
		d.nd[d.lastNode].up = r
		d.nd[d.lastNode].down = k

		if len(name) == 1 {
			d.nd[d.lastNode].color = 0
			d.nd[d.lastNode].scolor = ""
		} else if k >= d.second {
			c, _ := strconv.ParseInt(name[1], 36, 0)
			d.nd[d.lastNode].color = int(c)
			d.nd[d.lastNode].scolor = ":" + name[1]
		} else {
			return fmt.Errorf("primary item must be uncolored")
		}
	}

	if !pp {
		for d.lastNode > i {
			// Remove lastNode from its item list
			k := d.nd[d.lastNode].itm
			d.nd[k].itm--
			d.nd[k].color = i - 1
			q := d.nd[d.lastNode].up
			r := d.nd[d.lastNode].down
			d.nd[q].down = r
			d.nd[r].up = q
			d.lastNode--
		}
	} else {
		d.nd[i].down = d.lastNode
		d.lastNode++ // create the next spacer
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		d.options++
		d.nd[d.lastNode].up = i + 1
		d.nd[d.lastNode].itm = -d.options
	}

	return nil
}

func (d *DLX) InitialContent() {
	for i := 0; i < d.lastItm; i++ {
		fmt.Printf("%2d[%2s,%2d,%2d] ",
			i, d.cl[i].name, d.cl[i].prev, d.cl[i].next)
	}
	fmt.Println()
	for i := 0; i <= d.lastNode; i++ {
		if i%d.lastItm == 0 {
			fmt.Println()
		}
		fmt.Printf("%2d[%2d,%2d,%2d] ",
			i, d.nd[i].itm, d.nd[i].up, d.nd[i].down)
	}
	fmt.Println()
}
