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

var colorMap = map[string]color.RGBColor{
	"B": color.HEX("#000000", true),
	"T": color.HEX("#edd1d8", true),
	"U": color.HEX("#801dae", true),
	"V": color.HEX("#177cb0", true),
	"W": color.HEX("#0c8918", true),
	"X": color.HEX("#ff0097", true),
	"Y": color.HEX("#00bc12", true),
	"Z": color.HEX("#cccccc", true),
	"O": color.HEX("#b0a4e3", true),
	"P": color.HEX("#ffb3a7", true),
	"Q": color.HEX("#60281e", true),
	"R": color.HEX("#a1afc9", true),
	"S": color.HEX("#ffb61e", true),
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
