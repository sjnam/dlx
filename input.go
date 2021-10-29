package dlx

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

func (x *XCC) inputItemNames(line string) error {
	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if x.second != maxCols {
				return fmt.Errorf("item name line contains | twice")
			}
			x.second = x.lastItm
			continue
		}

		if strings.ContainsAny(itm, ":|") {
			return fmt.Errorf("illegal character in item name")
		}
		x.cl[x.lastItm].name = itm

		if len(itm) > 8 {
			return fmt.Errorf("item name too long")
		}

		// Check for duplicate item name
		var k int
		for k = 1; x.cl[k].name != x.cl[x.lastItm].name; k++ {
		}
		if k < x.lastItm {
			return fmt.Errorf("duplicate item name")
		}

		// Initialize lastItm to a new item with an empty list
		if x.lastItm > maxCols {
			return fmt.Errorf("too many items")
		}
		x.cl[x.lastItm-1].next = x.lastItm
		x.cl[x.lastItm].prev = x.lastItm - 1 // nd[lastItm].len = 0
		x.nd[x.lastItm].down = x.lastItm
		x.nd[x.lastItm].up = x.nd[x.lastItm].down
		x.lastItm++
	}

	if x.second == maxCols {
		x.second = x.lastItm
	}

	x.cl[x.lastItm].prev = x.lastItm - 1
	x.cl[x.lastItm-1].next = x.lastItm

	x.cl[x.second].prev = x.lastItm
	x.cl[x.lastItm].next = x.second

	x.cl[root].prev = x.second - 1
	x.cl[x.second-1].next = root
	x.lastNode = x.lastItm

	return nil
}

func (x *XCC) inputOptions(line string) error {
	var (
		pp = false
		i  = x.lastNode // remember the spacer at the left of this option
	)

	for _, opt := range strings.Fields(line) {
		if len(opt) > 8 {
			return fmt.Errorf("item name too long")
		}
		owc := strings.Split(opt, ":")
		x.cl[x.lastItm].name = owc[0]

		// Create a node for the item named in opt
		var k int
		for k = 0; x.cl[k].name != x.cl[x.lastItm].name; k++ {
		}
		if k == x.lastItm {
			return fmt.Errorf("unknown item name")
		}
		if x.nd[k].color >= i {
			return fmt.Errorf("duplicate item name in this option")
		}
		x.lastNode++
		if x.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		x.nd[x.lastNode].itm = k
		if k < x.second {
			pp = true
		}

		// Insert node lastNode into the list item k
		t := x.nd[k].itm + 1
		x.nd[k].itm = t
		x.nd[k].color = x.lastNode

		t = rand.Intn(t)
		var r int
		for r = k; t != 0; r, t = x.nd[r].down, t-1 {
		}
		q := x.nd[r].up
		x.nd[r].up = x.lastNode
		x.nd[q].down = x.nd[r].up
		x.nd[x.lastNode].up = q
		x.nd[x.lastNode].down = r

		if len(owc) == 1 {
			x.nd[x.lastNode].color = 0
		} else if k >= x.second {
			if len(owc[1]) != 1 {
				return fmt.Errorf("color must be a single character")
			}
			x.nd[x.lastNode].color = int(owc[1][0])
		} else {
			return fmt.Errorf("primary item must be uncolored")
		}
	}

	if !pp {
		for x.lastNode > i {
			k := x.nd[x.lastNode].itm
			x.nd[k].itm--
			x.nd[k].color = i - 1
			q := x.nd[x.lastNode].up
			r := x.nd[x.lastNode].down
			x.nd[q].down = r
			x.nd[r].up = q
			x.lastNode--
		}
	} else {
		x.nd[i].down = x.lastNode
		x.lastNode++ // create the next spacer
		if x.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		x.options++
		x.nd[x.lastNode].up = i + 1
		x.nd[x.lastNode].itm = -x.options
	}

	return nil
}

func (x *XCC) InputMatrix(rd io.Reader) error {
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
		x.lastItm = 1

		if err := x.inputItemNames(line); err != nil {
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

		if err := x.inputOptions(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (x *XCC) InitialContent() {
	for i := 0; i <= x.lastItm; i++ {
		fmt.Printf("(%s,%d,%d) ", x.cl[i].name, x.cl[i].prev, x.cl[i].next)
	}
	fmt.Println()
	for i := 0; i <= x.lastNode; i++ {
		if i%7 == 0 {
			fmt.Println()
		}
		fmt.Printf("(%d,%2d,%2d,%c) ",
			x.nd[i].itm, x.nd[i].up, x.nd[i].down, byte(x.nd[i].color))
	}
	fmt.Println()
}
