package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sjnam/dlx"
)

func spinner(delay time.Duration) {
	s := "."
	for i := 1; ; i++ {
		if i%40 == 0 {
			fmt.Printf("\n")
			s = "."
			i = 1
		}
		fmt.Printf("\r%s", s)
		time.Sleep(delay)
		s = s + "."
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
	if len(args) != 2 {
		log.Fatalf("%s dlx-file\n", args[0])
	}

	dlxInput := args[1]
	fd, err := os.Open(dlxInput)
	if err != nil {
		log.Fatal(err)
	}

	dlxInput = filepath.Base(dlxInput)
	name := strings.Split(dlxInput, ".")
	dimen := strings.Split(name[0], "x")
	nr, _ := strconv.Atoi(dimen[0])
	nc, _ := strconv.Atoi(dimen[1])

	go spinner(100 * time.Millisecond)

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(), fd)
	if err != nil {
		log.Fatal(err)
	}

	box := make([][]int, nr)
	for i := range box {
		box[i] = make([]int, nc)
	}

	solution := <-solStream
	for _, opt := range solution {
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

	fmt.Printf("\n")
	for j := 0; j < nr; j++ {
		for k := 0; k < nc; k++ {
			fmt.Printf("%d ", box[j][k])
		}
		fmt.Println()
	}
}
