# Covering with multiplicities and colors via Dancing Links
Go implementation of Donald Knuth's Algorithm 7.2.2.1M
for covering with multiplicities and colors. This implementation is based on the
Algorithm M described in https://www-cs-faculty.stanford.edu/~knuth/fasc5c.ps.gz

## Usage
````go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sjnam/dlx"
)

func main() {
	input := `
| A simple example of color controls
A B 2:3|C | X Y
A B X:0 Y:0
A C X:1 Y:1
C X:0
B X:1
C Y:1
`
	mcc := dlx.NewMCC()
	mcc.Debug = true
	res := mcc.Dance(strings.NewReader(input))

	for sol := range res.Solutions {
		for _, opt := range sol {
			fmt.Println(opt)
		}
	}
}

// (5 options, 3+2 items, 19 entries successfully read)
// [B X:1]
// [A C X:1 Y:1]
// [C Y:1]
// [null C]
// Altogether 1 solution 19 updates, 4 cleansings, 6 nodes.
````

## Examples

### Langford pairing
````
$ cd examples/langford
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
$ go run main.go puzzles.txt
Q[    1]: ..43..2.9..5..9..1.7..6..43..6..2.8719...74...5..83...6.....1.5..35.869..4291.3..
A[    1]: 864371259325849761971265843436192587198657432257483916689734125713528694542916378
Q[    2]: .4.1...5.1.7..396.52...8..........17...9.68..8.3.5.62..9..6.5436...8.7..25..971..
A[    2]: 346179258187523964529648371965832417472916835813754629798261543631485792254397186
Q[    3]: 6..12.384..8459.72.....6..5...264.3..7..8...694...3...31.....5..897.....5.2...19.
A[    3]: 695127384138459672724836915851264739273981546946573821317692458489715263562348197
Q[    4]: 4972.....1..4....5....16.9862.3...4.3..9.......1.726....2..587....6....453..97.61
A[    4]: 497258316186439725253716498629381547375964182841572639962145873718623954534897261
Q[    5]: ..591.3.8..94.3.6..275..1...3....2.1...82...7..6..7..4....8....64.15.7..89....42.
A[    5]: 465912378189473562327568149738645291954821637216397854573284916642159783891736425
...
Q[70098]: ..2.....9.3...25....61..37..........2..4..13...7..6.4...18.....76...54....9..76..
A[70098]: 472653819138792564956148372694531287285479136317286945521864793763915428849327651
Q[70099]: .3............1..87..58........24.5..4.8739....36.....9.......2..5..2.912.....7.4
A[70099]: 438297165659431278721586349167924853542873916893615427974168532385742691216359784
Solving took: 6.535064875s
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
What is Word search? https://thewordsearch.com/
````
$ cd examples/wordsearch
$ go run main.go movie.txt 13 13
봄게하대위게하밀은마기생충
돼여힘다장시제국설국열차타
지변름의인여의변해밀양며신
가명호가도봄날은간다리생과
우량장인을원엽가더날파활함
물생화파괴겨강기휘마이의께
에인홍물밤골울기적올란발죄
빠한련과막것극그친인드견와
진콤낮동의태아구리씨그보벌
날달투나택시운전사고가녀이
다컴는억추의인살박쥐봄아카
웰수다오수정씨자금한절친바
복스시아오하행산부적의공공

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
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟦🟦🟦🟦🟦🟦🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟦🟦🟦🟦🟦🟦🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟦🟦🟦🟦🟦🟦🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟦🟦🟦🟦🟦🟦🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟦🟦🟦🟦🟦🟦🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟦🟦🟦🟦🟦🟦🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟫🟫🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟫🟫🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟦🟦🟦🟦🟦🟦🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟦🟦🟦🟦🟦🟦🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟦🟦🟦🟦🟦🟦🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟦🟦🟦🟦🟦🟦🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟦🟦🟦🟦🟦🟦🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟪🟪🟪🟪🟦🟦🟦🟦🟦🟦🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟧🟧🟧🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟧🟧🟧🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪🟧🟧🟧🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟪🟪🟪🟪⬛🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟩🟩🟩🟩🟩🟨🟨🟨🟨🟨🟨🟨
🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟩🟩🟩🟩🟩🟪🟪🟪🟪🟧🟧🟧
🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟩🟩🟩🟩🟩🟪🟪🟪🟪🟧🟧🟧
🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟩🟩🟩🟩🟩🟪🟪🟪🟪🟧🟧🟧
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟩🟩🟩🟩🟩🟪🟪🟪🟪🟧🟧🟧
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟫🟫🟧🟧🟧
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟫🟫🟧🟧🟧
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟨🟨🟨🟨🟨🟨🟨🟩🟩🟩🟩🟩
````
