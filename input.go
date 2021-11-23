package dlx

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const delta int = 4

func inputMatrix(m *MCC, rd io.Reader) error {
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
		m.lastItm = 1
		break
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := inputItemNames(m, line); err != nil {
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

		if err := inputOptions(m, line); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if m.Debug {
		fmt.Fprintf(os.Stderr,
			"(%d options, %d+%d items, %d entries successfully read)\n",
			m.options, m.second-1, m.lastItm-m.second, m.lastNode-m.lastItm)
	}

	return nil
}

func inputItemNames(m *MCC, line string) error {
	if m.lastItm == 0 {
		return fmt.Errorf("no items")
	}

	for _, itm := range strings.Fields(line) {
		if itm == "|" {
			if m.second != maxCols {
				return fmt.Errorf("item name line contains | twice")
			}
			m.second = m.lastItm
			continue
		}

		if len(itm) > maxNameLength {
			return fmt.Errorf("item name too long")
		}

		q, r := 1, 1
		m.cl[m.lastItm].name = itm
		if strings.Contains(itm, "|") {
			bn := strings.Split(itm, "|")
			m.cl[m.lastItm].name = bn[1]
			q, _ = strconv.Atoi(bn[0])
			r = q
			if strings.Contains(bn[0], ":") {
				bounds := strings.Split(bn[0], ":")
				r, _ = strconv.Atoi(bounds[0])
				q, _ = strconv.Atoi(bounds[1])
			}
		}
		m.cl[m.lastItm].bound = q
		m.cl[m.lastItm].slack = q - r

		// Check for duplicate item name
		for k := 1; k < m.lastItm; k++ {
			if m.cl[k].name == m.cl[m.lastItm].name {
				return fmt.Errorf("duplicate item name")
			}
		}

		// Initialize lastItm to a new item with an empty list
		if m.lastItm > maxCols {
			return fmt.Errorf("too many items")
		}

		m.cl[m.lastItm-1].next = m.lastItm
		m.cl[m.lastItm].prev = m.lastItm - 1

		m.nd[m.lastItm].down = m.lastItm
		m.nd[m.lastItm].up = m.lastItm
		m.lastItm++
		if m.lastItm >= len(m.cl)-delta {
			m.cl = append(m.cl, make([]item, chunkSize)...)
			m.nd = append(m.nd, make([]node, chunkSize)...)
		}
	}

	if m.second == maxCols {
		m.second = m.lastItm
	}

	m.cl[m.lastItm].prev = m.lastItm - 1
	m.cl[m.lastItm-1].next = m.lastItm

	m.cl[m.second].prev = m.lastItm
	m.cl[m.lastItm].next = m.second
	// this sequence works properly whether second == lastItm

	m.cl[root].prev = m.second - 1
	m.cl[m.second-1].next = root

	m.lastNode = m.lastItm // reserve all the header nodes and the first spacer

	return nil
}

func inputOptions(m *MCC, line string) error {
	leftSpacer := m.lastNode // remember the spacer at the left of this option
	nonePrimary := true
	for _, opt := range strings.Fields(line) {
		if len(opt) > maxNameLength {
			return fmt.Errorf("item name too long")
		}
		if opt[0] == ':' {
			return fmt.Errorf("empty item name")
		}
		m.cl[m.lastItm].name = opt
		icr := strings.Index(opt, ":")
		if icr >= 0 { // has color code
			m.cl[m.lastItm].name = opt[:icr]
		}

		// Create a node for the item named in opt
		k := 0
		for ; m.cl[k].name != m.cl[m.lastItm].name; k++ {
		}
		if k == m.lastItm {
			return fmt.Errorf("unknown item name")
		}
		if m.nd[k].color >= leftSpacer { // aux field
			return fmt.Errorf("duplicate item name")
		}

		m.lastNode++
		if m.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}

		if m.lastNode >= len(m.nd)-delta {
			m.nd = append(m.nd, make([]node, chunkSize)...)
		}

		m.nd[m.lastNode].itm = k
		if k < m.second {
			nonePrimary = false
		}

		t := m.nd[k].itm + 1
		m.nd[k].itm = t
		m.nd[k].color = m.lastNode // aux field
		r := k
		for t = rand.Intn(t); t > 0; r = m.nd[r].down {
			t--
		}
		q := m.nd[r].up
		m.nd[q].down = m.lastNode
		m.nd[r].up = m.lastNode
		m.nd[m.lastNode].up = q
		m.nd[m.lastNode].down = r

		m.nd[m.lastNode].color = 0
		m.nd[m.lastNode].colorName = ""
		if icr >= 0 { // has color code
			if k >= m.second {
				cName := opt[icr+1:]
				m.nd[m.lastNode].color = int(crc32.ChecksumIEEE([]byte(cName)))
				m.nd[m.lastNode].colorName = ":" + cName
			} else {
				return fmt.Errorf("primary item must be uncolored")
			}
		}
	}

	if nonePrimary { // Option ignored (no primary items)
		for m.lastNode > leftSpacer { // Remove lastNode from its item list
			k := m.nd[m.lastNode].itm
			m.nd[k].itm--
			m.nd[k].color = leftSpacer - 1
			q, r := m.nd[m.lastNode].up, m.nd[m.lastNode].down
			m.nd[q].down, m.nd[r].up = r, q
			m.lastNode--
		}
	} else {
		m.nd[leftSpacer].down = m.lastNode
		m.lastNode++ // create the next spacer
		if m.lastNode == maxNodes {
			return fmt.Errorf("too many nodes")
		}
		if m.lastNode >= len(m.nd)-delta {
			m.nd = append(m.nd, make([]node, chunkSize)...)
		}
		m.options++
		m.nd[m.lastNode].up = leftSpacer + 1
		m.nd[m.lastNode].itm = -m.options
	}

	return nil
}
