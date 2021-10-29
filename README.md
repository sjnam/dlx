# go-dlx
Go implementation of Donald Knuth's Algorithm 7.2.2.1C for exact cover with colors.

This code is based on the Algorithm C described in
http://www-cs-faculty.stanford.edu/~knuth/fasc5c.ps.gz

## example

### nqueen
````
$ go run examples/queen/main.go 4
1:
. Q . .
. . . Q
Q . . .
. . Q .

2:
. . Q .
Q . . .
. . . Q
. Q . .

````

### sudoku
````
$ cd examples/sudoku
$ go run main.go < s17.dlx
1:
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

### zebra puzzle
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
water       tea         milk        orange      coffee
yellow      blue        red         white       green
Norway      Ukraine     England     Spain       Japan
diplomat    nurse       sculptor    violinist   painter
fox         horse       snail       dog         zebra
````

### pentominoes
- 12 pieces: O P Q R S T U V W X Y Z

````
$ cd example/pentominoes
$ go run main.go 3 20
Solution: 1
U U X O O O O O Z W W T T T R Q Q Q Q V
U X X X P P Z Z Z Y W W T R R R S S Q V
U U X P P P Z Y Y Y Y W T R S S S V V V

Solution: 2
U U X O O O O O S S S R T W Y Y Y Y Z V
U X X X P P Q S S R R R T W W Y Z Z Z V
U U X P P P Q Q Q Q R T T T W W Z V V V

Solution: 3
U U X P P P Q Q Q Q R T T T W W Z V V V
U X X X P P Q S S R R R T W W Y Z Z Z V
U U X O O O O O S S S R T W Y Y Y Y Z V

Solution: 4
U U X P P P Z Y Y Y Y W T R S S S V V V
U X X X P P Z Z Z Y W W T R R R S S Q V
U U X O O O O O Z W W T T T R Q Q Q Q V

Solution: 5
V Z Y Y Y Y W T R S S S O O O O O X U U
V Z Z Z Y W W T R R R S S Q P P X X X U
V V V Z W W T T T R Q Q Q Q P P P X U U

Solution: 6
V Q Q Q Q R T T T W W Z O O O O O X U U
V Q S S R R R T W W Y Z Z Z P P X X X U
V V V S S S R T W Y Y Y Y Z P P P X U U

Solution: 7
V V V S S S R T W Y Y Y Y Z P P P X U U
V Q S S R R R T W W Y Z Z Z P P X X X U
V Q Q Q Q R T T T W W Z O O O O O X U U

Solution: 8
V V V Z W W T T T R Q Q Q Q P P P X U U
V Z Z Z Y W W T R R R S S Q P P X X X U
V Z Y Y Y Y W T R S S S O O O O O X U U
````
