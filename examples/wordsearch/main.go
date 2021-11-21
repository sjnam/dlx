package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

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

func wordSearchDLX(words []string, wd, ht int) io.Reader {
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
			wlen := utf8.RuneCountInString(word)
			runew := []rune(word)
			for _, a := range [][]rune{runew, reverse(runew)} {
				for r := 0; r < ht; r++ {
					for c := 0; c < wd; c++ {
						// horizontal placement
						if c+wlen <= wd {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < wlen; i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r, c+i, a[i])
							}
							fmt.Fprintln(pw)
						}
						// vertical placement
						if r+wlen <= ht {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < wlen; i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r+i, c, a[i])
							}
							fmt.Fprintln(pw)
						}
						//upward diagonal placement
						if r+wlen <= ht && c-wlen+1 >= 0 {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < wlen; i++ {
								fmt.Fprintf(pw, "%x%x:%c ", r+i, c-i, a[i])
							}
							fmt.Fprintln(pw)
						}
						// downward diagonal placement
						if r+wlen <= ht && c+wlen <= wd {
							fmt.Fprintf(pw, "%s ", word)
							for i := 0; i < wlen; i++ {
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

func puzzleBoard(sol []dlx.Option, wd, ht int) {
	board := make([][]rune, ht)
	for i := 0; i < ht; i++ {
		board[i] = make([]rune, wd)
	}

	for _, opt := range sol {
		opt = opt[1:]
		for _, pos := range opt {
			chr := strings.Split(pos, ":")
			x, y := decode(chr[0])
			board[x][y], _ = utf8.DecodeRuneInString(chr[1])
		}
	}

	ascii := true
	for i := 0; i < wd; i++ {
		if board[0][i] != 0 {
			ascii = utf8.RuneLen(board[0][i]) <= 1
			break
		}
	}

	for i := 0; i < ht; i++ {
		for j := 0; j < wd; j++ {
			if board[i][j] == 0 {
				if ascii {
					board[i][j] = rune('A' + rand.Intn(26))
				} else {
					board[i][j] = rune('가' + rand.Intn(int('힣'-'가')))
				}
			}
			fmt.Printf("%c", board[i][j])
			if ascii {
				fmt.Print(" ")
			}
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

	d := dlx.NewMCC()
	res := d.Dance(
		wordSearchDLX(strings.Fields(string(buf)), wd, ht))

	i := 0
	for sol := range res.Solutions {
		i++
		fmt.Printf("%d:\n", i)
		puzzleBoard(sol, wd, ht)
		fmt.Println()
	}
}
