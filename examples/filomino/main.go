package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/sjnam/dlx"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("usage: %s w d\n", args[0])
		return
	}

	nr, _ := strconv.Atoi(args[1])
	nc, _ := strconv.Atoi(args[2])

	d, err := dlx.NewDLX(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}

	box := make([][]int, nr)
	for i := range box {
		box[i] = make([]int, nc)
	}

	i := 0
	for sol := range d.Dance() {
		i++
		fmt.Printf("Solution: %d\n", i)
		for _, opt := range sol {
			sort.Strings(opt)
			n := 0
			for j := 0; j < len(opt); j++ {
				if len(opt[j]) == 2 {
					n++
				} else {
					break
				}
			}
			for j := 0; j < n; j++ {
				rc, _ := strconv.Atoi(opt[j])
				r := rc / 10
				c := rc - r*10
				box[r][c] = n
			}
		}

		for j := 0; j < nr; j++ {
			for k := 0; k < nc; k++ {
				fmt.Printf("%d ", box[j][k])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
