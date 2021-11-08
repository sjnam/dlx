package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sjnam/dlx"
)

func bx(j, k int) int {
	return j/3*3 + k/3
}

func sudokuDLX(rd io.Reader) io.Reader {
	var pos, row, col, box [9][9]int

	buf := make([]byte, 9)
	var c, j int
	for {
		_, err := io.ReadFull(rd, buf)
		if err == io.EOF {
			break
		}

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

func sudokuSolve(line string) <-chan [][]byte {
	sudokuStream := make(chan [][]byte)
	go func() {
		defer close(sudokuStream)

		r := strings.NewReader(strings.TrimSpace(line))
		var buff bytes.Buffer

		rd := io.TeeReader(r, &buff)
		d := dlx.NewDancer()

		solStream, err := d.Dance(context.Background(), sudokuDLX(rd))
		if err != nil {
			log.Fatal(err)
		}

		res := make([][]byte, 2)
		qu, _ := ioutil.ReadAll(&buff)
		res[0] = qu
		board := make([]byte, len(qu))
		copy(board, qu)

		solution := <-solStream
		for _, opt := range solution {
			var x, y int
			var z byte
			fin := 0
			for _, v := range opt {
				if v[0] == 'p' {
					x = int(v[1] - '0')
					y = int(v[2] - '0')
					fin++
				}
				if v[0] == 'r' {
					z = v[2]
					fin++
				}
				if fin == 2 {
					break
				}
			}

			board[x*9+y] = z
		}
		res[1] = board
		sudokuStream <- res
	}()

	return sudokuStream
}

func fanIn(
	channels ...<-chan [][]byte,
) <-chan [][]byte {
	var wg sync.WaitGroup
	multiplexedStream := make(chan [][]byte)

	multiplex := func(c <-chan [][]byte) {
		defer wg.Done()
		for i := range c {
			select {
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("usage: %s dlx-file\n", args[0])
	}

	fd, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	var solutions []<-chan [][]byte

	scnr := bufio.NewScanner(fd)
	for scnr.Scan() {
		solutions = append(solutions, sudokuSolve(scnr.Text()))
	}
	if err := scnr.Err(); err != nil {
		log.Fatal(err)
	}

	i := 0
	for sol := range fanIn(solutions...) {
		i++
		fmt.Printf("Q[%2d]: %s\n", i, string(sol[0]))
		fmt.Printf("A[%2d]: %s\n", i, string(sol[1]))
	}
}
