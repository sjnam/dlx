package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sjnam/dlx"
)

func sudokuDLX(rd io.Reader) io.Reader {
	var c, j int
	var pos, row, col, box [9][9]int
	buf := make([]byte, 9)

	for {
		if _, err := io.ReadFull(rd, buf); err == io.EOF {
			break
		}

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

func sudokuSolver(stream <-chan string) <-chan [][]byte {
	ch := make(chan [][]byte)

	go func() {
		defer close(ch)

		for line := range stream {
			d := dlx.NewMCC()
			res := d.Dance(sudokuDLX(strings.NewReader(line)))

			ans := []byte(line)
			for _, opt := range <-res.Solutions {
				x := int(opt[0][1] - '0')
				y := int(opt[0][2] - '0')
				ans[x*9+y] = opt[1][2]
			}

			ch <- [][]byte{[]byte(line), ans}
		}
	}()

	return ch
}

func fanIn(ctx context.Context, channels []<-chan [][]byte) <-chan [][]byte {
	var wg sync.WaitGroup
	multiplexedStream := make(chan [][]byte)

	multiplex := func(c <-chan [][]byte) {
		defer wg.Done()
		for s := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexedStream <- s:
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

	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fd, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	numSolvers := runtime.NumCPU()
	var genStream []chan string
	for i := 0; i < numSolvers; i++ {
		genStream = append(genStream, make(chan string))
	}

	go func() {
		defer func() {
			for _, gs := range genStream {
				close(gs)
			}
		}()

		i := 0
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			genStream[i%numSolvers] <- strings.TrimSpace(scanner.Text())
			i++
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	var solvers []<-chan [][]byte
	for _, gs := range genStream {
		solvers = append(solvers, sudokuSolver(gs))
	}

	i := 0
	for s := range fanIn(ctx, solvers) {
		i++
		fmt.Printf("Q[%5d]: %s\n", i, s[0])
		fmt.Printf("A[%5d]: %s\n", i, s[1])
	}

	fmt.Printf("Solving took: %v\n", time.Since(start))
}
