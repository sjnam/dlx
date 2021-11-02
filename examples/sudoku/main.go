package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/sjnam/dlx"
)

func bx(j, k int) int {
	return j/3*3 + k/3
}

func sudokuDLX(rd io.Reader) io.Reader {
	var pos, row, col, box [9][9]int

	var c, j int
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		buf := scanner.Text()
		for k := 0; k < 9; k++ {
			if buf[k] >= '1' && buf[k] <= '9' {
				d := int(buf[k] - '1')
				x := bx(j, k)
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
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	r, w := io.Pipe()
	go func() {
		defer func() {
			_ = w.Close()
		}()

		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if pos[j][k] == 0 {
					fmt.Fprintf(w, "p%d%d ", j, k)
				}
			}
		}
		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if row[j][k] == 0 {
					fmt.Fprintf(w, "r%d%d ", j, k+1)
				}
			}
		}
		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if col[j][k] == 0 {
					fmt.Fprintf(w, "c%d%d ", j, k+1)
				}
			}
		}
		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				if box[j][k] == 0 {
					fmt.Fprintf(w, "b%d%d ", j, k+1)
				}
			}
		}
		fmt.Fprintln(w)

		for j := 0; j < 9; j++ {
			for k := 0; k < 9; k++ {
				for d := 0; d < 9; d++ {
					x := bx(j, k)
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
	var buff bytes.Buffer

	rd := io.TeeReader(os.Stdin, &buff)
	d, err := dlx.NewDLX(sudokuDLX(rd))
	if err != nil {
		fmt.Println(err)
		return
	}

	var board [][]string
	scanner := bufio.NewScanner(&buff)
	for scanner.Scan() {
		buf := scanner.Text()
		board = append(board, strings.Split(buf, ""))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	i := 0
	for solution := range d.Dance() {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range solution {
			sort.Strings(opt)
			board[opt[2][1]-'0'][opt[2][2]-'0'] = string(opt[3][2])
		}
		for i, row := range board {
			for j, col := range row {
				fmt.Printf("%s ", col)
				if (j+1)%3 == 0 {
					fmt.Print(" ")
				}
			}
			fmt.Println()
			if (i+1)%3 == 0 && i != 8 {
				fmt.Println()
			}
		}
		fmt.Println()
	}
}
