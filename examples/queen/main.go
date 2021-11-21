package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/sjnam/dlx"
)

func encode(x int) byte {
	if x < 10 {
		return byte('0') + byte(x)
	} else if x < 36 {
		return byte('a') + byte(x) - 10
	} else {
		return byte('A') + byte(x) - 36
	}
}

func queenDLX(n int) io.Reader {
	nn := n + n - 2
	if nn > 62 {
		log.Fatal("Sorry , I can't currently handle n>32!")
	}

	r, w := io.Pipe()
	go func() {
		defer func() {
			_ = w.Close()
		}()

		for j := 0; j < n; j++ {
			t := n + j
			if j&1 != 0 {
				t = n - 1 - j
			}
			t = t >> 1
			fmt.Fprintf(w, "r%c c%c ", encode(t), encode(t))
		}
		fmt.Fprint(w, "|")
		for j := 1; j < nn; j++ {
			fmt.Fprintf(w, " a%c b%c", encode(j), encode(j))
		}
		fmt.Fprintln(w)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				fmt.Fprintf(w, "r%c c%c", encode(j), encode(k))
				t := j + k
				if t != 0 && t < nn {
					fmt.Fprintf(w, " a%c", encode(t))
				}
				t = n - 1 - j + k
				if t != 0 && t < nn {
					fmt.Fprintf(w, " b%c", encode(t))
				}
				fmt.Fprintln(w)
			}
		}
	}()

	return r
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("usage: %s n\n", args[0])
		return
	}
	n, _ := strconv.Atoi(os.Args[1])

	mcc := dlx.NewMCC()
	res := mcc.Dance(queenDLX(n))

	i := 0
	board := make([][]string, n)
	for r := 0; r < n; r++ {
		board[r] = make([]string, n)
	}
	for solution := range res.Solutions {
		i++
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				board[r][c] = "."
			}
		}
		var found int
		for _, opt := range solution {
			var r, c int64
			found = 0
			for _, rc := range opt {
				switch t := rc[0]; t {
				case 'r':
					r, _ = strconv.ParseInt(rc[1:], 32, 0)
					found++
				case 'c':
					c, _ = strconv.ParseInt(rc[1:], 32, 0)
					found++
				}
				if found == 2 {
					break
				}
			}
			board[r][c] = "Q"
		}

		fmt.Printf("%d:\n", i)
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				fmt.Printf("%s ", board[r][c])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
