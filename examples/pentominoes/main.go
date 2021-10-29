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

	fname := fmt.Sprintf("%dx%d.dlx", nr, nc)
	fp, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	xcc := dlx.NewXCC()
	err = xcc.InputMatrix(fp)
	if err != nil {
		fmt.Println(err)
		return
	}

	box := make([][]string, nr)
	for i := range box {
		box[i] = make([]string, nc)
	}

	i := 0
	for sol := range xcc.Dance() {
		i++
		fmt.Printf("Solution: %d\n", i)
		for _, opt := range sol {
			sort.Slice(opt, func(i, j int) bool {
				return opt[i] > opt[j]
			})

			c := opt[0]
			for j := 1; j < len(opt); j++ {
				x, _ := strconv.ParseInt(string(opt[j][0]), 36, 0)
				y, _ := strconv.ParseInt(string(opt[j][1]), 36, 0)
				box[x][y] = c
			}
		}

		for j := 0; j < nr; j++ {
			for k := 0; k < nc; k++ {
				c := " "
				if box[j][k] != "" {
					c = box[j][k]
				}
				fmt.Printf("%s ", c)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
