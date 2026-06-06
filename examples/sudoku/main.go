package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/sjnam/dlx"
	"github.com/sjnam/ofanin"
)

type puzzle [81]byte

func sudokuDLX(line puzzle) io.Reader {
	var pos, row, col, box [9][9]int
	for i, j := 0, 0; i < 81; i, j = i+9, j+1 {
		for k := 0; k < 9; k++ {
			if ch := line[i+k]; ch >= '1' && ch <= '9' {
				d := int(ch - '1')
				x := j/3*3 + k/3
				pos[j][k], row[j][d], col[k][d], box[x][d] = d+1, k+1, j+1, j+1
			}
		}
	}

	var b bytes.Buffer
	b.Grow(8192)
	// Primary items: every unfilled position/row/col/box constraint.
	for j := 0; j < 9; j++ {
		for k := 0; k < 9; k++ {
			if pos[j][k] == 0 {
				fmt.Fprintf(&b, "p%d%d ", j, k)
			}
		}
	}
	for _, t := range []struct {
		c string
		a *[9][9]int
	}{{"r", &row}, {"c", &col}, {"b", &box}} {
		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if t.a[j][k] == 0 {
					fmt.Fprintf(&b, "%s%d%d ", t.c, j, k+1)
				}
			}
		}
	}
	b.WriteByte('\n')
	// One option per legal (position, digit) placement.
	for j := 0; j < 9; j++ {
		for k := 0; k < 9; k++ {
			for d := 0; d < 9; d++ {
				x := j/3*3 + k/3
				if pos[j][k] == 0 && row[j][d] == 0 && col[k][d] == 0 && box[x][d] == 0 {
					fmt.Fprintf(&b, "p%d%d r%d%d c%d%d b%d%d\n", j, k, j, d+1, k, d+1, x, d+1)
				}
			}
		}
	}
	return bytes.NewReader(b.Bytes())
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("usage: %s file\n", args[0])
	}

	start := time.Now()
	defer func() {
		fmt.Printf("Solving took: %v\n", time.Since(start))
	}()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	dlxSudoku := ofanin.NewOrderedFanIn[puzzle, [2]puzzle](ctx)
	dlxSudoku.InputStream = func() <-chan puzzle {
		ch := make(chan puzzle)
		go func() {
			defer close(ch)
			data, err := os.ReadFile(args[1])
			if err != nil {
				panic(err)
			}
			for len(data) > 0 {
				line, rest, _ := bytes.Cut(data, []byte{'\n'})
				ch <- puzzle(line)
				data = rest
			}
		}()
		return ch
	}()
	dlxSudoku.DoWork = func(q puzzle) [2]puzzle {
		xc := dlx.NewDancer()
		a := q
		res := xc.Dance(sudokuDLX(q))
		for _, opt := range <-res.Solutions {
			x := int(opt[0][1] - '0')
			y := int(opt[0][2] - '0')
			a[x*9+y] = byte(opt[1][2])
		}
		return [2]puzzle{q, a}
	}

	i := 0
	for s := range dlxSudoku.Process() {
		i++
		fmt.Printf("Q[%5d]: %s\n", i, s[0])
		fmt.Printf("A[%5d]: %s\n", i, s[1])
	}
}
