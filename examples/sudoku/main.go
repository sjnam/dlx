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
	var c, j int
	var pos, row, col, box [9][9]int

	for i := 0; i < 81; i += 9 {
		buf := line[i : i+9]
		for k := 0; k < 9; k++ {
			if buf[k] >= '1' && buf[k] <= '9' {
				d := int(buf[k] - '1')
				x := j/3*3 + k/3
				pos[j][k] = d + 1
				if row[j][d] != 0 {
					log.Fatalf("digit %d appears in columns"+
						" %d and %d of row %d!", d+1, row[j][d]-1, k, j)
				}
				row[j][d] = k + 1
				if col[k][d] != 0 {
					log.Fatalf("digit %d appears in rows"+
						" %d and %d of column %d!", d+1, col[k][d]-1, j, k)
				}
				col[k][d] = j + 1
				if box[x][d] != 0 {
					log.Fatalf("digit %d appears in rows"+
						" %d and %d of box %d!", d+1, box[x][d]-1, j, x)
				}
				box[x][d] = j + 1
				c++
			}
		}
		j++
	}

	r, w := io.Pipe()
	go func() {
		defer func() {
			_ = w.Close()
		}()

		for j = 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if pos[j][k] == 0 {
					fmt.Fprintf(w, "p%d%d ", j, k)
				}
			}
		}
		for j = 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if row[j][k] == 0 {
					fmt.Fprintf(w, "r%d%d ", j, k+1)
				}
			}
		}
		for j = 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if col[j][k] == 0 {
					fmt.Fprintf(w, "c%d%d ", j, k+1)
				}
			}
		}
		for j = 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if box[j][k] == 0 {
					fmt.Fprintf(w, "b%d%d ", j, k+1)
				}
			}
		}
		fmt.Fprintln(w)

		for j = 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				for d := 0; d < 9; d++ {
					x := j/3*3 + k/3
					if pos[j][k] == 0 && row[j][d] == 0 &&
						col[k][d] == 0 && box[x][d] == 0 {
						fmt.Fprintf(w, "p%d%d r%d%d c%d%d b%d%d\n",
							j, k, j, d+1, k, d+1, x, d+1)
					}
				}
			}
		}
	}()

	return r
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
