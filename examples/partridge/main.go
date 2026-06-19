package main

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/sjnam/dlx"
)

func patridgeDLX(n int) io.Reader {
	N := n * (n + 1) / 2
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		for i := 1; i <= n; i++ {
			fmt.Fprintf(w, "%d:%d|#%d ", i, i, i)
		}
		for i := 0; i < N; i++ {
			for j := 0; j < N; j++ {
				fmt.Fprintf(w, "%d,%d ", i, j)
			}
		}
		fmt.Fprintln(w)
		for t := 1; t <= n; t++ {
			for r := 0; r < N-t+1; r++ {
				for c := 0; c < N-t+1; c++ {
					fmt.Fprintf(w, "#%d ", t)
					for rr := 0; rr < t; rr++ {
						for cc := 0; cc < t; cc++ {
							fmt.Fprintf(w, "%d,%d ", r+rr, c+cc)
						}
					}
					fmt.Fprintln(w)
				}
			}
		}
	}()
	return r
}

// fillBoard fills size[r][c] with tile size and tile[r][c] with a unique tile ID.
func fillBoard(n int, sol []dlx.Option) (size, tile [][]int) {
	N := n * (n + 1) / 2
	size = make([][]int, N)
	tile = make([][]int, N)
	for i := range size {
		size[i] = make([]int, N)
		tile[i] = make([]int, N)
	}
	tileID := 0
	for _, opt := range sol {
		s, _ := strconv.Atoi(opt[0][1:]) // "#k" -> k
		tileID++
		for _, coord := range opt[1:] {
			parts := strings.Split(coord, ",")
			r, _ := strconv.Atoi(parts[0])
			c, _ := strconv.Atoi(parts[1])
			size[r][c] = s
			tile[r][c] = tileID
		}
	}
	return
}

// printBoard draws the tiling as a Unicode box-drawing diagram.
//
// To stay compact while keeping each tile roughly square, every base cell maps
// to a single text row (1 tall) and one interior column (1 wide), so the board
// is 2N+1 chars wide and N+1 rows tall. Each text row merges the horizontal grid
// line at the top of a cell row with that row's interior, rather than spending a
// separate row on it. With the terminal's ~2:1 character height/width ratio, a
// k×k tile renders as (2k-1)×k characters, which displays approximately square.
func printBoard(size, tile [][]int) {
	N := len(tile)

	tileAt := func(r, c int) int {
		if r < 0 || r >= N || c < 0 || c >= N {
			return 0
		}
		return tile[r][c]
	}

	// hBorder: horizontal line above board row r at column c.
	hBorder := func(r, c int) bool {
		return r == 0 || r == N || tileAt(r, c) != tileAt(r-1, c)
	}

	// vBorder: vertical line left of board column c at row r.
	vBorder := func(r, c int) bool {
		return c == 0 || c == N || tileAt(r, c) != tileAt(r, c-1)
	}

	// Label each tile in its bottom-right cell, where the cell's top edge is
	// interior to the tile (no horizontal line to collide with). A 1×1 tile has
	// no such cell — its only row is also its top border — so it stays blank.
	type pos struct{ r, c int }
	tileCells := make(map[int][]pos)
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			tileCells[tile[r][c]] = append(tileCells[tile[r][c]], pos{r, c})
		}
	}
	label := make([][]int, N)
	for i := range label {
		label[i] = make([]int, N)
	}
	for _, cells := range tileCells {
		maxR, maxC := cells[0].r, cells[0].c
		for _, p := range cells {
			if p.r > maxR {
				maxR = p.r
			}
			if p.c > maxC {
				maxC = p.c
			}
		}
		if s := size[cells[0].r][cells[0].c]; s > 1 {
			label[maxR][maxC] = s
		}
	}

	junc := [16]rune{
		' ', '╶', '╴', '─',
		'╷', '┌', '┐', '┬',
		'╵', '└', '┘', '┴',
		'│', '├', '┤', '┼',
	}

	var sb strings.Builder
	for r := 0; r <= N; r++ {
		sb.Reset()
		for c := 0; c <= N; c++ {
			// Junction at the top-left corner of cell (r, c): vertical edges run
			// up (row r-1) and down (row r); horizontal edges run left and right.
			b := 0
			if r > 0 && vBorder(r-1, c) {
				b |= 8
			}
			if r < N && vBorder(r, c) {
				b |= 4
			}
			if c > 0 && hBorder(r, c-1) {
				b |= 2
			}
			if c < N && hBorder(r, c) {
				b |= 1
			}
			sb.WriteRune(junc[b])

			if c < N {
				switch {
				case hBorder(r, c):
					sb.WriteRune('─')
				case label[r][c] > 0:
					sb.WriteRune(rune('0' + label[r][c]))
				default:
					sb.WriteRune(' ')
				}
			}
		}
		fmt.Println(sb.String())
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Minute)
	defer cancel()

	n := 8
	mcc := dlx.NewDancer()
	mcc.PulseInterval = 30 * time.Second
	mcc = mcc.WithContext(ctx)
	res := mcc.Dance(patridgeDLX(n))

	start := time.Now()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case st, ok := <-res.Heartbeat:
				if !ok {
					return
				}
				fmt.Printf("%s (%v)\n", st, time.Since(start))
			}
		}
	}()

	i := 0
	for sol := range res.Solutions {
		i++
		fmt.Printf("%d:\n", i)
		size, tile := fillBoard(n, sol)
		printBoard(size, tile)
	}
}
