package dlx

import (
	"fmt"
	"strings"
)

func ExampleXCC() {
	input := `
|A simple example of color controls
A B C | X Y
A B X:0 Y:0
A C X:1 Y:1
X:0 Y:1
B X:1
C Y:1
`
	d, err := NewDLX(strings.NewReader(input))
	if err != nil {
		fmt.Println(err)
		return
	}

	for sol := range d.Dance() {
		for _, opt := range sol {
			fmt.Println(opt)
		}
	}

	// Unordered output:
	// [A C X:1 Y:1]
	// [B X:1]
}

func ExampleDLX() {
	input := `
| A simple example
A B C D E | F G
C E F
A D G
B C F
A D
B G
D E G
`
	d, err := NewDLX(strings.NewReader(input))
	if err != nil {
		fmt.Println(err)
		return
	}

	for sol := range d.Dance() {
		for _, opt := range sol {
			fmt.Println(opt)
		}
	}

	// Unordered output:
	// [C E F]
	// [B G]
	// [A D]
}
