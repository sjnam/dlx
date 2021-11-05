package main

import (
	"fmt"
	"strings"

	"github.com/sjnam/go-dlx"
)

const dlxInput = `
|Zebra Puzzle
#1 #2 #3 #4 #5 #6 #7 #8 #9 #10 #11 #12 #13 #14 #15 #16 | N0 N1 N2 N3 N4 J0 J1 J2 J3 J4 P0 P1 P2 P3 P4 D0 D1 D2 D3 D4 C0 C1 C2 C3 C4
|The Englishman lives in a red house.
#1 N0:England C0:red
#1 N1:England C1:red
#1 N2:England C2:red
#1 N3:England C3:red
#1 N4:England C4:red
|The painter comes from Japan.
#2 N0:Japan J0:painter
#2 N1:Japan J1:painter
#2 N2:Japan J2:painter
#2 N3:Japan J3:painter
#2 N4:Japan J4:painter
|The yellow house hosts a diplomat.
#3 J0:diplomat C0:yellow
#3 J1:diplomat C1:yellow
#3 J2:diplomat C2:yellow
#3 J3:diplomat C3:yellow
#3 J4:diplomat C4:yellow
|The coffee-lover's house is green.
#4 D0:coffee C0:green
#4 D1:coffee C1:green
#4 D2:coffee C2:green
#4 D3:coffee C3:green
#4 D4:coffee C4:green
|The Norwegian's house is hte leftmost.
#5 N0:Norway
|The dog's owner is from Spain.
#6 N0:Spain P0:dog
#6 N1:Spain P1:dog
#6 N2:Spain P2:dog
#6 N3:Spain P3:dog
#6 N4:Spain P4:dog
|The milk drinker lives in the middle house.
#7 D2:milk
|The violinist drinks orange juice.
#8 J0:violinist D0:orange
#8 J1:violinist D1:orange
#8 J2:violinist D2:orange
#8 J3:violinist D3:orange
#8 J4:violinist D4:orange
|The white house is just left of the green one.
#9 C0:white C1:green
#9 C1:white C2:green
#9 C2:white C3:green
#9 C3:white C4:green
|The Ukrainian drinks tea.
#10 N0:Ukraine D0:tea
#10 N1:Ukraine D1:tea
#10 N2:Ukraine D2:tea
#10 N3:Ukraine D3:tea
#10 N4:Ukraine D4:tea
|The Norwegian lives next to the blue house.
#11 N0:Norway C1:blue
#11 N1:Norway C2:blue
#11 N2:Norway C3:blue
#11 N3:Norway C4:blue
#11 C0:blue N1:Norway
#11 C1:blue N2:Norway
#11 C2:blue N3:Norway
#11 C3:blue N4:Norway
|The sculptor breeds snails.
#12 J0:sculptor P0:snail
#12 J1:sculptor P1:snail
#12 J2:sculptor P2:snail
#12 J3:sculptor P3:snail
#12 J4:sculptor P4:snail
|The horse lives next to the diplomat.
#13 J0:diplomat P1:horse
#13 J1:diplomat P2:horse
#13 J2:diplomat P3:horse
#13 J3:diplomat P4:horse
#13 P0:horse J1:diplomat
#13 P1:horse J2:diplomat
#13 P2:horse J3:diplomat
#13 P3:horse J4:diplomat
|The nurse lives next to the fox.
#14 J0:nurse P1:fox
#14 J1:nurse P2:fox
#14 J2:nurse P3:fox
#14 J3:nurse P4:fox
#14 P0:fox J1:nurse
#14 P1:fox J2:nurse
#14 P2:fox J3:nurse
#14 P3:fox J4:nurse
|Somebody trains a zebra.
#15 P0:zebra
#15 P1:zebra
#15 P2:zebra
#15 P3:zebra
#15 P4:zebra
|Somebody prefers to drink just plain water
#16 D0:water
#16 D1:water
#16 D2:water
#16 D3:water
#16 D4:water
`

func main() {
	d, err := dlx.NewDLX(strings.NewReader(dlxInput))
	if err != nil {
		fmt.Println(err)
		return
	}

	answer := map[byte][]string{
		'N': make([]string, 5),
		'J': make([]string, 5),
		'P': make([]string, 5),
		'D': make([]string, 5),
		'C': make([]string, 5),
	}

	solution := <-d.Dance()
	for _, opt := range solution {
		for i := 1; i < len(opt); i++ {
			kv := strings.Split(opt[i], ":")
			answer[kv[0][0]][kv[0][1]-'0'] = kv[1]
		}
	}

	for _, line := range answer {
		for _, v := range line {
			fmt.Printf("%-12s", v)
		}
		fmt.Println()
	}
}
