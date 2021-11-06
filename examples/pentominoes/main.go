package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sjnam/dlx"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("usage: %s dlx-file\n", args[0])
	}

	dlxInput := args[1]
	name := strings.Split(dlxInput, ".")
	dimen := strings.Split(name[0], "x")
	nr, _ := strconv.Atoi(dimen[0])
	nc, _ := strconv.Atoi(dimen[1])

	fd, err := os.Open(dlxInput)
	if err != nil {
		log.Fatal(err)
	}

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(), fd)
	if err != nil {
		log.Fatal(err)
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
