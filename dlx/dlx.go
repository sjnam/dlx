package dlx

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"
)

const (
	infty    = int(^uint(0) >> 1)
	maxCount = infty
	root     = 0             // cl[root] is the gateway to the unsettled items
	maxLevel = 500           // at most this many options in a solution
	maxCols  = 10000         // at most this many items
	maxNodes = 100000000     // at most this many nonzero elements in the matrix
	maxLine  = 9*maxCols + 3 // a size big enough to hold all item names

	maxNameLength = 32 // max item name length
)

var (
	ErrInputLineTooLong    = errors.New("input line way too long")
	ErrNoItems             = errors.New("no items")
	ErrEmptyItemName       = errors.New("empty item name")
	ErrUnknownItemName     = errors.New("unknown item name")
	ErrItemNameTooLong     = errors.New("item name too long")
	ErrDuplicateItemName   = errors.New("duplicate item name")
	ErrIllegalCharacter    = errors.New("illegal character in item name")
	ErrIllegalItemNameLine = errors.New("item name line contains | twice")
	ErrTooManyItems        = errors.New("too many items")
	ErrTooManyNodes        = errors.New("too many nodes")
	ErrTooManyLevels       = errors.New("too many levels")
	ErrPrimaryItemColored  = errors.New("primary item must be uncolored")
)

type node struct {
	up, down  int    // predecessor and successor in item list
	itm       int    // the item containing this node
	color     int    // the color specified by this node, if any
	colorName string // the color name string
}

type item struct {
	name         string // symbolic identification of the item, for printing
	prev, next   int    // neighbors of this item
	bound, slack int    // residual capacity of ths item
}

// DLX dancing links object
type DLX struct {
	nd       []node // the master list of nodes
	lastNode int    // the first node in nd that's not yet used
	cl       []item // the master list of items
	lastItm  int    // the first item in cl that's not yet used
	second   int    // boundary between primary and secondary items
}

type XC struct {
	*DLX
}

type XCC struct {
	*DLX
}

type MCC struct {
	*DLX
}

func newDLX() *DLX {
	return &DLX{
		nd:     make([]node, maxNodes),
		cl:     make([]item, maxCols+2),
		second: maxCols,
	}
}

// NewXC generates a dancing machine.
func NewXC() XC {
	return XC{
		DLX: newDLX(),
	}
}

// NewXCC generates a dancing machine.
func NewXCC() XCC {
	return XCC{
		DLX: newDLX(),
	}
}

// NewMCC generates a dancing machine.
func NewMCC() MCC {
	return MCC{
		DLX: newDLX(),
	}
}

// Dancer solves exact cover problem while dancing.
type Dancer interface {
	InputItemNames(string) error
	InputOptions(string) error
	Dance(context.Context, io.Reader) (<-chan [][]string, error)
}

func inputMatrix(m Dancer, rd io.Reader) error {
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
			line = ""
			continue
		}
		break
	}

	if err := m.InputItemNames(line); err != nil {
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

		if err := m.InputOptions(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
