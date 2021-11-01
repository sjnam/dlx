package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sjnam/dlx"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func digit(b byte) int {
	r := b
	if r >= 'a' && r <= 'f' {
		r = r - 'a' + 10
	} else {
		r = r - '0'
	}

	return int(r)
}

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("usage: %s w d\n", args[0])
		return
	}

	nr, _ := strconv.Atoi(args[1])
	nc, _ := strconv.Atoi(args[2])
	if nr > 16 || nc > 16 {
		fmt.Println("w and d should be <= 16")
		return
	}

	go spinner(100 * time.Millisecond)

	d, err := dlx.NewDLX(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}

	box := make([][]int, nr)
	for i := range box {
		box[i] = make([]int, nc)
	}

	sol := <-d.Dance()
	for _, opt := range sol {
		n := 0
		var coor [][2]int
		for _, c := range opt {
			if len(c) == 2 {
				n++
				coor = append(coor, [2]int{digit(c[0]), digit(c[1])})
			}
		}
		for i := 0; i < n; i++ {
			box[coor[i][0]][coor[i][1]] = n
		}
	}

	fmt.Printf("\r")
	for j := 0; j < nr; j++ {
		for k := 0; k < nc; k++ {
			fmt.Printf("%d ", box[j][k])
		}
		fmt.Println()
	}
}
