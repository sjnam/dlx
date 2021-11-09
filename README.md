# Dancing Links
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

### Langford pairing
````
$ go run main.go 4
[2 3 4 2 1 3 1 4]
````

### Pentominoes
- 12 pieces: **O P Q R S T U V W X Y Z**

````
$ cd example/pentominoes
$ go run main.go 8x8.dlx
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
$ cd examples/sudoku
$ go run main.go top96.dlx 
Q[ 1]: 9.4..5...25.6..1..31......8.7...9...4..26......147....7.......2...3..8.6.4.....9.
A[ 1]: 964815237258637149317924658872159364495263781631478925783596412529341876146782593
Q[ 2]: ..247..58..............1.4.....2...9528.9.4....9...1.........3.3....75..685..2...
A[ 2]: 132479658847563291956281347413725869528196473769348125271854936394617582685932714
Q[ 3]: .476...5.8.3.....2.....9......8.5..6...1.....6.24......78...51...6....4..9...4..7
A[ 3]: 947628351863751492125349678734895126589162734612473985478236519256917843391584267
Q[ 4]: ..5..8..18......9.......78....4.....64....9......53..2.6.........138..5....9.714.
A[ 4]: 935748621876231594124695783512469378643872915789153462267514839491386257358927146
...
Q[96]: 4.....8.5.3..........7......2.....6.....5.4......1.......6.3.7.5..2.....1.9......
A[96]: 417369825638125947952748316825437169791856432346912758284693571573281694169574283
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

$ go run main.go 10x10.filomino.dlx 
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
````


### Word Search
소설 토지의 목차들을 찾아라.
````
$ cd examples/wordsearch
$ go run main.go land.txt 13 13
하바다마아발과말종태동기삶
역병과흉년혼백의귀향타하의
카령아차바닥모를늪속으로형
다악후이세만막사의희명다태
꿈속의귀마동사카들들차라우
아타마카다파나매람떠파빛풍
세월을넘고나은사나혈속하의
파어하사하젊는는고으바다국
모두긴여로하자과로아하타북
음운나아일남결울서과촌정용
과계타에는순하어둠의발소리
적절밤자운명적인것사바하라
추바사아카리자의실인래거귀

$ go run main.go mathematicians.txt 15 15
H I L B E R T L C A N T O R O 
X O L F W Q F F O H H C R I K 
G T C A T A L A N R E T S D O 
C T Q S M I N K O W S K I R T 
G E S Y L V E S T E R V R A M 
F N A W N O R R E P Y K T M E 
S I A B E L U A D N A L V A L 
J T W E I E R S T R A S S D L 
S U I Z T I W R U H L X C A I 
U D Z E C R E H S I A L G H N 
B O R E L W D N A R T R E B X 
D G R A M T S U I N E B O R F 
R U N G E O J J E N S E N C X 
F F O K R A M E E T I M R E H 
E K N O P P L E S N E H X S N 
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

### Partridge puzzle
````
$ go run example/partridge/main.go
1:
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩⬜🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟥🟥🟥🟥🟥🟥🟥🟥🟦🟦🟦🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟪🟪⬛⬛🟦🟦🟦🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟪🟪⬛⬛🟦🟦🟦🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟪🟪🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟪🟪🟪🟪🟪🟪🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟪🟪🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟧🟧🟧🟧⬛⬛🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟪🟪🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧⬛⬛🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟪🟪🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟦🟦🟦🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟦🟦🟦🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟦🟦🟦🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟦🟦🟦🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟦🟦🟦🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
🟧🟧🟧🟧🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟪🟦🟦🟦🟩🟩🟩🟩🟩🟧🟧🟧🟧🟥🟥🟥🟥🟥🟥🟥🟥
````
