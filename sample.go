package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/sjnam/go-dlx/dlx"
)

func main() {
	dlxInput := `
| A simple example
A B C D E | F G
C E F
A D G
B C F
A D
B G
D E G
`
	dx, err := dlx.NewDancer(strings.NewReader(dlxInput))
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for sol := range dx.Dance() {
		i++
		fmt.Printf("%d:\n", i)
		for _, opt := range sol {
			fmt.Println(opt)
		}
		fmt.Println()
	}
}
