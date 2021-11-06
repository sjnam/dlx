package main

import (
	"context"
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

	fp, err := os.Open(fmt.Sprintf("%dx%d.dlx", nr, nc))
	if err != nil {
		fmt.Println(err)
		return
	}

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(), fp)
	if err != nil {
		fmt.Println(err)
		return
	}

	box := make([][]string, nr)
	for i := range box {
		box[i] = make([]string, nc)
		for j := range box[i] {
			box[i][j] = " "
		}
	}

	i := 0
	for solution := range solStream {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range solution {
			sort.Strings(opt)
			c := opt[len(opt)-1]
			for j := 0; j < len(opt)-1; j++ {
				x, _ := strconv.ParseInt(string(opt[j][0]), 36, 0)
				y, _ := strconv.ParseInt(string(opt[j][1]), 36, 0)
				box[x][y] = c
			}
		}

		for j := 0; j < nr; j++ {
			for k := 0; k < nc; k++ {
				fmt.Printf("%s ", box[j][k])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
