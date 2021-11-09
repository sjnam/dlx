package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sjnam/dlx"
)

func reverse(sa []rune) []rune {
	r := make([]rune, len(sa))
	copy(r, sa)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}

func wordSearch(words []string, wd, ht int) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer func() {
			_ = pw.Close()
		}()

		// item line
		fmt.Fprint(pw, strings.Join(words, " "))
		fmt.Fprint(pw, " | ")
		for i := 0; i < ht; i++ {
			for j := 0; j < wd; j++ {
				fmt.Fprintf(pw, "%x%x ", i, j)
			}
		}
		fmt.Fprintln(pw)

		// option lines
		for _, word := range words {
			nm := []rune(word)
			for _, a := range [][]rune{nm, reverse(nm)} {
				for r := 0; r < ht; r++ {
					for c := 0; c < wd; c++ {
						// horizontal placement
						if c+len(word) <= wd {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < len(word); i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r, c+i, a[i])
							}
							fmt.Fprintln(pw)
						}
						// vertical placement
						if r+len(word) <= ht {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < len(word); i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r+i, c, a[i])
							}
							fmt.Fprintln(pw)
						}
						//upward diagonal placement
						if r+len(word) <= ht && c-len(word)+1 >= 0 {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < len(word); i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r+i, c-i, a[i])
							}
							fmt.Fprintln(pw)
						}
						// downward diagonal placement
						if r+len(word) <= ht && c+len(word) <= wd {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < len(word); i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r+i, c+i, a[i])
							}
							fmt.Fprintln(pw)
						}
					}
				}
			}
		}
	}()

	return pr
}

func decode(str string) (int, int) {
	var x [2]int
	for i := 0; i < 2; i++ {
		if str[i] < '0'+10 {
			x[i] = int(str[i] - '0')
		} else if str[i] >= 'a' {
			x[i] = int(str[i] - 'a' + 10)
		}
	}
	return x[0], x[1]
}

func puzzleBoard(sol [][]string, wd, ht int) {
	board := make([][]string, ht)
	for i := 0; i < ht; i++ {
		board[i] = make([]string, wd)
	}

	for _, opt := range sol {
		opt = opt[1:]
		for _, pos := range opt {
			chr := strings.Split(pos, ":")
			x, y := decode(chr[0])
			board[x][y] = chr[1]
		}
	}

	for i := 0; i < ht; i++ {
		for j := 0; j < wd; j++ {
			if board[i][j] == "" {
				board[i][j] = "_"
			}
			fmt.Printf("%s ", board[i][j])
		}
		fmt.Println()
	}
}

func main() {
	args := os.Args
	if len(args) != 4 {
		log.Fatalf("usage: %s input-file wd ht\n", args[0])
	}

	wd, _ := strconv.Atoi(args[2])
	ht, _ := strconv.Atoi(args[3])
	buf, err := ioutil.ReadFile(args[1])
	if err != nil {
		log.Fatal(err)
	}

	d := dlx.NewDancer()
	solStream, err := d.Dance(context.Background(),
		wordSearch(strings.Fields(string(buf)), wd, ht))
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for sol := range solStream {
		i++
		fmt.Printf("%d:\n", i)
		puzzleBoard(sol, wd, ht)
		fmt.Println()
	}
}
