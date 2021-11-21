package dlx

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"strings"
)

const delta int = 4

func inputMatrix(d *Dancer, rd io.Reader) error {
	var line string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > maxLine {
			return fmt.Errorf("input line way too long")
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' { // bypass comment or blank line
			continue
		}
		d.lastItm = 1
		break
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := inputItemNames(d, line); err != nil {
		return err
	}

	for scanner.Scan() {
		line = scanner.Text()
		if len(line) > maxLine {
			return fmt.Errorf("input line way too long")
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '|' { // bypass comment or blank line
			continue
		}

		if err := inputOptions(d, line); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if d.Debug {
		fmt.Fprintf(os.Stderr,
			"(%d options, %d+%d items, %d entries successfully read)\n",
			d.options, d.second-1, d.lastItm-d.second, d.lastNode-d.lastItm)
	}

	return nil
}

func inputItemNames(d *Dancer, line string) error {
	if d.lastItm == 0 {
		return fmt.Errorf("no items")
	}

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
		d.cl[d.lastItm].name = itm
		if strings.Contains(itm, "|") {
			bn := strings.Split(itm, "|")
			d.cl[d.lastItm].name = bn[1]
			q, _ = strconv.Atoi(bn[0])
			r = q
			if strings.Contains(bn[0], ":") {
				bounds := strings.Split(bn[0], ":")
				r, _ = strconv.Atoi(bounds[0])
				q, _ = strconv.Atoi(bounds[1])
			}
		}
		d.cl[d.lastItm].bound = q
		d.cl[d.lastItm].slack = q - r

		// Check for duplicate item name
		for k := 1; k < d.lastItm; k++ {
			if d.cl[k].name == d.cl[d.lastItm].name {
				return fmt.Errorf("duplicate item name")
			}
		}

		// Initialize lastItm to a new item with an empty list
		if d.lastItm > maxCols {
			return fmt.Errorf("too many items")
		}

		d.cl[d.lastItm-1].next = d.lastItm
		d.cl[d.lastItm].prev = d.lastItm - 1

		d.nd[d.lastItm].down = d.lastItm
		d.nd[d.lastItm].up = d.lastItm
		d.lastItm++
		if d.lastItm >= len(d.cl)-delta {
			d.cl = append(d.cl, make([]item, chunkSize)...)
			d.nd = append(d.nd, make([]node, chunkSize)...)
		}
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

	return nil
}

func inputOptions(d *Dancer, line string) error {
	leftSpacer := d.lastNode // remember the spacer at the left of this option
	nonePrimary := true
	for _, opt := range strings.Fields(line) {
		if len(opt) > maxNameLength {
			return fmt.Errorf("item name too long")
		}
		if opt[0] == ':' {
			return fmt.Errorf("empty item name")
		}
		d.cl[d.lastItm].name = opt
		icr := strings.Index(opt, ":")
		if icr >= 0 { // has color code
			d.cl[d.lastItm].name = opt[:icr]
		}

		// Create a node for the item named in opt
		k := 0
		for ; d.cl[k].name != d.cl[d.lastItm].name; k++ {
		}
		if k == d.lastItm {
			return fmt.Errorf("unknown item name")
		}
		if d.nd[k].color >= leftSpacer { // aux field
			return fmt.Errorf("duplicate item name")
		}

		d.lastNode++
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}

		if d.lastNode >= len(d.nd)-delta {
			d.nd = append(d.nd, make([]node, chunkSize)...)
		}

		d.nd[d.lastNode].itm = k
		if k < d.second {
			nonePrimary = false
		}

		d.nd[k].itm++
		d.nd[k].color = d.lastNode // aux field
		r := d.nd[k].up
		d.nd[r].down = d.lastNode
		d.nd[k].up = d.lastNode
		d.nd[d.lastNode].up = r
		d.nd[d.lastNode].down = k

		d.nd[d.lastNode].color = 0
		d.nd[d.lastNode].colorName = ""
		if icr >= 0 { // has color code
			if k >= d.second {
				cName := opt[icr+1:]
				d.nd[d.lastNode].color = int(crc32.ChecksumIEEE([]byte(cName)))
				d.nd[d.lastNode].colorName = ":" + cName
			} else {
				return fmt.Errorf("primary item must be uncolored")
			}
		}
	}

	if nonePrimary { // Option ignored (no primary items)
		for d.lastNode > leftSpacer { // Remove lastNode from its item list
			k := d.nd[d.lastNode].itm
			d.nd[k].itm--
			d.nd[k].color = leftSpacer - 1
			q, r := d.nd[d.lastNode].up, d.nd[d.lastNode].down
			d.nd[q].down, d.nd[r].up = r, q
			d.lastNode--
		}
	} else {
		d.nd[leftSpacer].down = d.lastNode
		d.lastNode++ // create the next spacer
		if d.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		if d.lastNode >= len(d.nd)-delta {
			d.nd = append(d.nd, make([]node, chunkSize)...)
		}
		d.options++
		d.nd[d.lastNode].up = leftSpacer + 1
		d.nd[d.lastNode].itm = -d.options
	}

	return nil
}
