package dlx

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func solve(matrix string) {
	d := NewDancer()
	d.Info = true
	solStream, err := d.Dance(context.Background(), strings.NewReader(matrix))
	if err != nil {
		fmt.Println(err)
		return
	}

	for sol := range solStream {
		for _, opt := range sol {
			sort.Strings(opt)
			fmt.Println(opt)
		}
	}

	d.Statistics()
}

func ExampleDancer_xc() {
	xcInput := `

| A simple example

A B C D E | F G
C E F
A D G


B C F
A D
B G
D E G
`
	solve(xcInput)

	// Unordered Output:
	// [A D]
	// [B G]
	// [C E F]
}

func ExampleDancer_xcc() {
	xccInput := `
|A simple example of color controls
A B C | X Y
A B X:0 Y:0
A C X:1 Y:1
| comments
X:0 Y:1
B X:1
C Y:1
`
	solve(xccInput)

	// Unordered Output:
	// [A C X:1 Y:1]
	// [B X:1]
}

func ExampleDancer_mcc() {
	mccInput := `
| A simple example of color controls

| and multiplicity

A B 2:3|C | X Y
A B X:0 Y:0

A C X:1 Y:1
C X:0

B X:1
C Y:1
| finished
`
	solve(mccInput)

	// Unordered Output:
	// [A C X:1 Y:1]
	// [B X:1]
	// [C Y:1]
	// [null C]
}
