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
// Cell width W=3 compensates for the ~2:1 terminal character height/width ratio,
// making each tile appear approximately square.
func printBoard(size, tile [][]int) {
	N := len(tile)
	const W = 3 // chars per cell interior; 3 wide × 1 tall ≈ square with 2:1 font ratio

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

	// Determine label position (center cell) for each tile.
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
		minR, minC := cells[0].r, cells[0].c
		for _, p := range cells {
			if p.r < minR {
				minR = p.r
			}
			if p.c < minC {
				minC = p.c
			}
		}
		s := size[cells[0].r][cells[0].c]
		label[minR+(s-1)/2][minC+(s-1)/2] = s
	}

	junc := [16]rune{
		' ', '╶', '╴', '─',
		'╷', '┌', '┐', '┬',
		'╵', '└', '┘', '┴',
		'│', '├', '┤', '┼',
	}

	var sb strings.Builder
	for dr := 0; dr <= 2*N; dr++ {
		sb.Reset()
		if dr%2 == 0 {
			// Border row: junctions and W-char horizontal segments.
			r := dr / 2
			for c := 0; c <= N; c++ {
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
					seg := ' '
					if hBorder(r, c) {
						seg = '─'
					}
					for range W {
						sb.WriteRune(seg)
					}
				}
			}
		} else {
			// Content row: vertical borders and W-char cell interiors.
			r := (dr - 1) / 2
			for c := 0; c <= N; c++ {
				if vBorder(r, c) {
					sb.WriteRune('│')
				} else {
					sb.WriteRune(' ')
				}
				if c < N {
					if lbl := label[r][c]; lbl > 0 {
						// Center the digit within W chars.
						pad := (W - 1) / 2
						for range pad {
							sb.WriteRune(' ')
						}
						sb.WriteRune(rune('0' + lbl))
						for range W - pad - 1 {
							sb.WriteRune(' ')
						}
					} else {
						for range W {
							sb.WriteRune(' ')
						}
					}
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
