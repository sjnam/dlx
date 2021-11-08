package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sjnam/dlx"
)

var color = [9]rune{
	'\U00002B1C', // white
	'\U00002B1B', // black
	'\U0001F7EB', // brown
	'\U0001F7E7', // orange
	'\U0001F7EA', // purple
	'\U0001F7E9', // green
	'\U0001F7E6', // blue
	'\U0001F7E8', // yellow
	'\U0001F7E5', // red
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func patridgeDLX(n int) io.Reader {
	N := n * (n + 1) / 2
	r, w := io.Pipe()
	go func() {
		defer func() {
			_ = w.Close()
		}()

		for i := 1; i <= n; i++ {
			fmt.Fprintf(w, "%d:%d|#%d ", i, i, i)
		}
		for i := 0; i < N; i++ {
			for j := 0; j < N; j++ {
				fmt.Fprintf(w, "%d,%d ", i, j)
			}
		}
		fmt.Fprintln(w)
		for t := 1; t <= n; t++ {
			for r := 0; r < N-t+1; r++ {
				for c := 0; c < N-t+1; c++ {
					fmt.Fprintf(w, "#%d ", t)
					for rr := 0; rr < t; rr++ {
						for cc := 0; cc < t; cc++ {
							fmt.Fprintf(w, "%d,%d ", r+rr, c+cc)
						}
					}
					fmt.Fprintln(w)
				}
			}
		}
	}()
	return r
}

func fillBoard(sol [][]string, board [][]rune) {
	for _, opt := range sol {
		sort.Strings(opt)
		s, _ := strconv.Atoi(opt[0][1:])
		for _, coord := range opt[1:] {
			co := strings.Split(coord, ",")
			r, _ := strconv.Atoi(co[0])
			c, _ := strconv.Atoi(co[1])
			board[r][c] = color[s]
		}
	}

	N := len(board)
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("%c", board[i][j])
		}
		fmt.Println()
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := 8
	d := dlx.NewDancer()
	solStream, err := d.Dance(ctx, patridgeDLX(n))
	if err != nil {
		log.Fatal(err)
	}

	go spinner(100 * time.Millisecond)

	N := n * (n + 1) / 2
	board := make([][]rune, N)
	for i := 0; i < len(board); i++ {
		board[i] = make([]rune, N)
	}

	i := 0
	for sol := range solStream {
		i++
		if i == 4 {
			cancel()
		}
		fmt.Printf("\n%d:\n", i)
		fillBoard(sol, board)
	}
}
