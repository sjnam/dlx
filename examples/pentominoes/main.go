package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/sjnam/dlx"
)

var colorMap = map[string]*color.Style256{
	"B": color.S256(0, 0),
	"T": color.S256(0, 225),
	"U": color.S256(0, 57),
	"V": color.S256(0, 27),
	"W": color.S256(0, 22),
	"X": color.S256(0, 198),
	"Y": color.S256(0, 48),
	"Z": color.S256(0, 253),
	"O": color.S256(0, 104),
	"P": color.S256(0, 172),
	"Q": color.S256(0, 94),
	"R": color.S256(0, 14),
	"S": color.S256(0, 11),
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("usage: %s dlx-file\n", args[0])
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

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(), fd)
	if err != nil {
		log.Fatal(err)
	}

	box := make([][]string, nr)
	for i := range box {
		box[i] = make([]string, nc)
		for j := range box[i] {
			box[i][j] = "B"
		}
	}

	i := 0
	for solution := range solStream {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range solution {
			for _, v := range opt[1:] {
				x, _ := strconv.ParseInt(string(v[0]), 36, 0)
				y, _ := strconv.ParseInt(string(v[1]), 36, 0)
				box[x][y] = opt[0]
			}
		}

		for j := 0; j < nr; j++ {
			for k := 0; k < nc; k++ {
				colorMap[box[j][k]].Print("  ")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
