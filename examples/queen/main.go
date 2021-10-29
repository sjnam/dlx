package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/sjnam/dlx"
)

func enocde(x int) byte {
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
		defer w.Close()

		for j := 0; j < n; j++ {
			t := n + j
			if j&1 != 0 {
				t = n - 1 - j
			}
			t = t >> 1
			fmt.Fprintf(w, "r%c c%c ", enocde(t), enocde(t))
		}
		fmt.Fprint(w, "|")
		for j := 1; j < nn; j++ {
			fmt.Fprintf(w, " a%c b%c", enocde(j), enocde(j))
		}
		fmt.Fprintln(w)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				fmt.Fprintf(w, "r%c c%c", enocde(j), enocde(k))
				t := j + k
				if t != 0 && t < nn {
					fmt.Fprintf(w, " a%c", enocde(t))
				}
				t = n - 1 - j + k
				if t != 0 && t < nn {
					fmt.Fprintf(w, " b%c", enocde(t))
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
	sn := os.Args[1]
	n, _ := strconv.Atoi(sn)

	xcc := dlx.NewXCC()
	err := xcc.InputMatrix(queenDLX(n))
	if err != nil {
		fmt.Println(err)
		return
	}

	i := 0
	for sol := range xcc.Dance() {
		i++
		brd := make([][]string, n)
		for r := 0; r < n; r++ {
			brd[r] = make([]string, n)
			for c := 0; c < n; c++ {
				brd[r][c] = "."
			}
		}
		for _, opt := range sol {
			var r, c int64
			for _, rc := range opt {
				if rc[0] == 'r' {
					r, _ = strconv.ParseInt(rc[1:], 32, 0)
				} else if rc[0] == 'c' {
					c, _ = strconv.ParseInt(rc[1:], 32, 0)
				}
			}
			brd[r][c] = "Q"
		}

		fmt.Printf("%d:\n", i)
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				fmt.Printf("%s ", brd[r][c])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
