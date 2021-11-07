package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sjnam/dlx"
)

var color = [9]rune{
	'\U0001F7EB', // brown
	'\U00002B1C', // white
	'\U00002B1B', // black
	'\U0001F7E7', // orange
	'\U0001F7EA', // purple
	'\U0001F7E9', // green
	'\U0001F7E6', // blue
	'\U0001F7E8', // yellow
	'\U0001F7E5', // red
}

func spinner(delay time.Duration) {
	s := "."
	for i := 1; ; i++ {
		if i%100 == 0 {
			fmt.Printf("\n")
			s = "."
			i = 1
		}
		fmt.Printf("\r%s", s)
		time.Sleep(delay)
		s = s + "."
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

		fmt.Fprintf(w, "\n")

		for t := 1; t <= n; t++ {
			for r := 0; r < N-t+1; r++ {
				for c := 0; c < N-t+1; c++ {
					fmt.Fprintf(w, "#%d ", t)
					for rr := 0; rr < t; rr++ {
						for cc := 0; cc < t; cc++ {
							fmt.Fprintf(w, "%d,%d ", r+rr, c+cc)
						}
					}
					fmt.Fprintf(w, "\n")
				}
			}
		}
	}()

	return r
}

func drawSquare(sol, board [][]string) {
	for _, opt := range sol {
		sort.Strings(opt)
		s, _ := strconv.Atoi(opt[0][1:])
		for _, coord := range opt[1:] {
			co := strings.Split(coord, ",")
			r, _ := strconv.Atoi(co[0])
			c, _ := strconv.Atoi(co[1])
			board[r][c] = fmt.Sprintf("%c", color[s])
		}
	}

	N := len(board)
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Print(board[i][j])
		}
		fmt.Println()
	}
}

func main() {
	go spinner(100 * time.Millisecond)

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("usage: %s n\n", args[0])
		return
	}
	n, _ := strconv.Atoi(os.Args[1])

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(), patridgeDLX(n))
	if err != nil {
		log.Fatal(err)
	}

	N := n * (n + 1) / 2
	board := make([][]string, N)
	for i := 0; i < len(board); i++ {
		board[i] = make([]string, N)
	}

	i := 0
	for sol := range solStream {
		i++
		fmt.Printf("%d:\n", i)
		drawSquare(sol, board)
	}
}
