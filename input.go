package dlx

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

const MaxNameLength = 32

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
		d.cl[d.lastItm].name = itm

		if len(itm) > MaxNameLength {
			return fmt.Errorf("item name too long")
		}

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
		d.cl[d.lastItm].prev = d.lastItm - 1 // nd[lastItm].len = 0
		d.nd[d.lastItm].down = d.lastItm
		d.nd[d.lastItm].up = d.nd[d.lastItm].down
		d.lastItm++
	}

	if d.second == maxCols {
		d.second = d.lastItm
	}

	d.cl[d.lastItm].prev = d.lastItm - 1
	d.cl[d.lastItm-1].next = d.lastItm

	d.cl[d.second].prev = d.lastItm
	d.cl[d.lastItm].next = d.second

	d.cl[root].prev = d.second - 1
	d.cl[d.second-1].next = root
	d.lastNode = d.lastItm

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
		owc := strings.Split(opt, ":")
		d.cl[d.lastItm].name = owc[0]

		// Create a node for the item named in opt
		var k int
		for k = 0; d.cl[k].name != d.cl[d.lastItm].name; k++ {
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
		d.nd[k].itm = t
		d.nd[k].color = d.lastNode

		t = rand.Intn(t)
		var r int
		for r = k; t != 0; r, t = d.nd[r].down, t-1 {
		}
		q := d.nd[r].up
		d.nd[r].up = d.lastNode
		d.nd[q].down = d.nd[r].up
		d.nd[d.lastNode].up = q
		d.nd[d.lastNode].down = r

		if len(owc) == 1 {
			d.nd[d.lastNode].color = 0
			d.nd[d.lastNode].scolor = ""
		} else if k >= d.second {
			c, _ := strconv.ParseInt(owc[1], 36, 0)
			d.nd[d.lastNode].color = int(c)
			d.nd[d.lastNode].scolor = ":" + owc[1]
		} else {
			return fmt.Errorf("primary item must be uncolored")
		}
	}

	if !pp {
		for d.lastNode > i {
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

func (d *DLX) inputMatrix(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > maxLine {
			return fmt.Errorf("input line way too long")
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' {
			continue
		}
		d.lastItm = 1

		if err := d.inputItemNames(line); err != nil {
			return err
		}
		break
	}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > maxLine {
			return fmt.Errorf("option line too long")
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' {
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

func (d *DLX) initialContent() {
	for i := 0; i <= d.lastItm; i++ {
		fmt.Printf("(%s,%d,%d) ", d.cl[i].name, d.cl[i].prev, d.cl[i].next)
	}
	fmt.Println()
	for i := 0; i <= d.lastNode; i++ {
		if i%7 == 0 {
			fmt.Println()
		}
		fmt.Printf("(%d,%2d,%2d,%c) ",
			d.nd[i].itm, d.nd[i].up, d.nd[i].down, byte(d.nd[i].color))
	}
	fmt.Println()
}
