# go-dlx
Go implementation of Donald Knuth's Algorithm 7.2.2.1M
for covering with multiplicities and colors. This implementation is based on the
Algorithm M described in https://www-cs-faculty.stanford.edu/~knuth/fasc5c.ps.gz

## Usage
````go
package dlx

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func solve(matrix string) {
	d := NewDancer()
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
}

func ExampleDancer_dlx() {
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
A B 2:3|C | X Y
A B X:0 Y:0
A C X:1 Y:1
C X:0
B X:1
C Y:1
`
	solve(mccInput)

	// Unordered Output:
	// [A C X:1 Y:1]
	// [B X:1]
	// [C Y:1]
	// [null C]
}
````

## Examples

### Nqueen
````
$ go run examples/queen/main.go 8
1:
Q . . . . . . . 
. . . . . Q . . 
. . . . . . . Q 
. . Q . . . . . 
. . . . . . Q . 
. . . Q . . . . 
. Q . . . . . . 
. . . . Q . . . 

2:
Q . . . . . . . 
. . . . Q . . . 
. . . . . . . Q 
. . . . . Q . . 
. . Q . . . . . 
. . . . . . Q . 
. Q . . . . . . 
. . . Q . . . . 

...
````

### Sudoku
````
| ......3..
| 1..4.....
| ......1.5
| 9........
| .....26..
| ....53...
| .5.8.....
| ...9...7.
| .83....4.

$ cd examples/sudoku
$ go run main.go < s17.dlx
5 9 7  2 1 8  3 6 4
1 3 2  4 6 5  8 9 7
8 6 4  3 7 9  1 2 5

9 1 5  6 8 4  7 3 2
3 4 8  7 9 2  6 5 1
2 7 6  1 5 3  4 8 9

6 5 9  8 4 7  2 1 3
4 2 1  9 3 6  5 7 8
7 8 3  5 2 1  9 4 6
````

### Zebra puzzle
Five people, from five different countries, have five different occupations,
own five different pets, drink five different beverages, and live in a row of
five different colored houses.

- The Englishman lives in a red house.
- The painter comes from Japan.
- The yellow house hosts a diplomat.
- The coffee-lover's house is green.
- The Norwegian's house is hte leftmost.
- The dog's owner is from Spain.
- The milk drinker lives in the middle house.
- The violinist drinks orange juice.
- The white house is just left of the green one.
- The Ukrainian drinks tea.
- The Norwegian lives next to the blue house.
- The sculptor breeds snails.
- The horse lives next to the diplomat.
- The nurse lives next to the fox.

Who trains the zebra, and who prefers to drink just plain water?

````
$ go run examples/zebra/main.go
Norway      Ukraine     England     Spain       Japan       
diplomat    nurse       sculptor    violinist   painter     
fox         horse       snail       dog         zebra       
water       tea         milk        orange      coffee      
yellow      blue        red         white       green
````

### Pentominoes
- 12 pieces: **O P Q R S T U V W X Y Z**

````
$ cd example/pentominoes
$ go run main.go 8 8
Solution: 1
P P P W U U U Y 
P P W W U T U Y 
Z W W T T T Y Y 
Z Z Z     T X Y 
R R Z     X X X 
V R R S S S X Q 
V R S S Q Q Q Q 
V V V O O O O O 

Solution: 2
W P P P U U U Y 
W W P P U T U Y 
Z W W T T T Y Y 
Z Z Z     T X Y 
R R Z     X X X 
V R R S S S X Q 
V R S S Q Q Q Q 
V V V O O O O O 

Solution: 3
V V V Z W W Q Q 
V Z Z Z R W W Q 
V Z P R R R W Q 
O P P     R X Q 
O P P     X X X 
O U U S S S X T 
O U S S Y T T T 
O U U Y Y Y Y T 

...
````

### Filomino
````
$ cd examples/filomino

| ..3.3...3.
| ..131...43
| 64...141..
| .6...4.4..
| 64...141..
| ..434...12
| ..3.3...2.
| ..434...12
| 24...161..
| .2...6.6..

$ go run main.go 10 10 < 10x10.filomino.dlx 
3 3 3 1 3 3 2 2 3 3 
4 4 1 3 1 3 4 4 4 3 
6 4 4 3 3 1 4 1 2 2 
6 6 6 6 4 4 1 4 4 4 
6 4 3 3 4 1 4 1 4 2 
4 4 4 3 4 3 4 4 1 2 
3 3 3 1 3 3 4 2 2 1 
2 4 4 3 4 4 6 6 1 2 
2 4 4 3 4 1 6 1 3 2 
1 2 2 3 4 6 6 6 3 3 

| ..4.....
| .32..26.
| .5....24
| ...46...
| ...33...
| 15....7.
| .21..24.
| .....1..

$ go run main.go 8 8 < 8x8.filomino.dlx
4 4 4 4 6 6 6 6 
3 3 2 6 2 2 6 6 
3 5 2 6 6 6 2 4 
1 5 4 4 6 6 2 4 
5 5 4 3 3 3 4 4 
1 5 4 7 7 7 7 1 
2 2 1 7 2 2 4 4 
3 3 3 7 7 1 4 4 
````
