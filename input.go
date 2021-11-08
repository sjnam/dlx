package dlx

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func (d *Dancer) inputMatrix(rd io.Reader) error {
	var line string

	scanner := bufio.NewScanner(rd)
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
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := d.inputItemNames(line); err != nil {
		return err
	}

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

		if err := d.inputOptions(line); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if d.Info {
		fmt.Fprintf(os.Stderr,
			"(%d options, %d+%d items, %d entries successfully read)\n",
			d.options, d.second-1, d.lastItm-d.second, d.lastNode-d.lastItm)
	}

	return nil
}

func (d *Dancer) inputItemNames(line string) error {
	if d.lastItm == 0 {
		return fmt.Errorf("no items")
	}

	cl, nd := d.cl, d.nd

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if d.second != maxCols {
				return fmt.Errorf("item name line contains | twice")
			}
			d.second = d.lastItm
			continue
		}

		if len(itm) > maxNameLength {
			return fmt.Errorf("item name too long")
		}

		q, r := 1, 1
		bn := strings.Split(itm, "|")
		if len(bn) == 1 {
			cl[d.lastItm].name = itm
		} else {
			cl[d.lastItm].name = bn[1]
			bounds := strings.Split(bn[0], ":")
			if len(bounds) == 1 {
				q, _ = strconv.Atoi(bounds[0])
				r = q
			} else {
				r, _ = strconv.Atoi(bounds[0])
				q, _ = strconv.Atoi(bounds[1])
			}
		}
		cl[d.lastItm].bound = q
		cl[d.lastItm].slack = q - r

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

func (d *Dancer) inputOptions(line string) error {
	var (
		cl = d.cl
		nd = d.nd
		i  = d.lastNode // remember the spacer at the left of this option
		pp = false
	)

	for _, opt := range strings.Fields(line) {
		if len(opt) > maxNameLength {
			return fmt.Errorf("item name too long")
		}
		if opt[0] == ':' {
			return fmt.Errorf("empty item name")
		}
		name := strings.Split(opt, ":")
		cl[d.lastItm].name = name[0]

		// Create a node for the item named in opt
		k := 0
		for ; cl[k].name != cl[d.lastItm].name; k++ {
		}
		if k == d.lastItm {
			return fmt.Errorf("unknown item name")
		}
		if nd[k].color >= i { // aux field
			return fmt.Errorf("duplicate item name")
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
		t := nd[k].itm + 1
		// we want to put the node into a random position of the list
		// we store the position of the new node into nd[k].color,
		// so that the test for duplicate items above will be correct.
		nd[k].itm = t            // len field; store the new length of the list
		nd[k].color = d.lastNode // aux field
		r := k
		for t = rand.Intn(t); t > 0; t-- {
			r = nd[r].down
		}
		q := nd[r].up
		nd[q].down = d.lastNode
		nd[r].up = d.lastNode
		nd[d.lastNode].up = q
		nd[d.lastNode].down = r

		nd[d.lastNode].color = 0
		nd[d.lastNode].colorName = ""
		if len(name) == 2 {
			if k >= d.second {
				c, _ := strconv.ParseInt(name[1], 36, 0)
				nd[d.lastNode].color = int(c)
				nd[d.lastNode].colorName = ":" + name[1]
			} else {
				return fmt.Errorf("primary item must be uncolored")
			}
		}
	}

	if !pp {
		for d.lastNode > i {
			// Remove lastNode from its item list
			k := nd[d.lastNode].itm
			nd[k].itm--
			nd[k].color = i - 1
			q, r := nd[d.lastNode].up, nd[d.lastNode].down
			nd[q].down, nd[r].up = r, q
			d.lastNode--
		}
	} else {
		nd[i].down = d.lastNode
		d.lastNode++ // create the next spacer
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		d.options++
		nd[d.lastNode].up = i + 1
		nd[d.lastNode].itm = -d.options
	}

	return nil
}
