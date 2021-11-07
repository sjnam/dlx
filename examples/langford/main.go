package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/sjnam/dlx"
)

func langfordDLX(n int) io.Reader {
	r, w := io.Pipe()

	go func() {
		defer func() {
			_ = w.Close()
		}()

		for i := 1; i <= n; i++ {
			fmt.Fprintf(w, "%d ", i)
		}
		for i := 1; i <= 2*n; i++ {
			fmt.Fprintf(w, "s%d ", i)
		}
		fmt.Fprintln(w)
		for i := 1; i <= n; i++ {
			for j := 1; j < 2*n-i; j++ {
				if i != n-(1-n&0x1) || j <= n/2 {
					fmt.Fprintf(w, "%d s%d s%d\n", i, j, j+i+1)
				}
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

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(), langfordDLX(n))
	if err != nil {
		log.Fatal(err)
	}

	s := make([]int, 2*n)
	for sol := range solStream {
		for _, opt := range sol {
			sort.Strings(opt)
			k, _ := strconv.Atoi(opt[0])
			for j := 1; j <= 2; j++ {
				p, _ := strconv.Atoi(opt[j][1:])
				s[p-1] = k
			}
		}
		fmt.Println(s)
	}
}
