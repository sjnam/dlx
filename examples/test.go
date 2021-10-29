package main

import (
	"fmt"
	"strings"

	"github.com/sjnam/dlx"
)

const dlxInput = `
|A simple example of color controls
A B C | X Y
A B X:0 Y:0
A C X:1 Y:1
X:0 Y:1
B X:1
C Y:1
`

func main() {
	xcc := dlx.NewXCC()
	err := xcc.InputMatrix(strings.NewReader(dlxInput))
	if err != nil {
		fmt.Println(err)
		return
	}

	for sol := range xcc.Dance() {
		for _, opt := range sol {
			fmt.Println(opt)
		}
	}
}
