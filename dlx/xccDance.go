package dlx

import (
	"context"
	"fmt"
	"io"
	"log"
)

func (xcc *XCC) getOption(p int) []string {
	var (
		option []string
		cl, nd = xcc.cl, xcc.nd
	)
	for q := p; ; {
		option = append(option,
			fmt.Sprintf("%s%s", cl[nd[q].itm].name, nd[q].colorName))
		q++
		if nd[q].itm <= 0 {
			q = nd[q].up
		}
		if q == p {
			break
		}
	}
	return option
}

func (xcc *XCC) cover(c int) {
	cl, nd := xcc.cl, xcc.nd
	l, r := cl[c].prev, cl[c].next
	cl[l].next = r
	cl[r].prev = l
	for rr := nd[c].down; rr >= xcc.lastItm; rr = nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if nd[nn].color >= 0 {
				uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				nd[uu].down = dd
				nd[dd].up = uu
				nd[cc].itm--
			}
			nn++
		}
	}
}

func (xcc *XCC) uncover(c int) {
	cl, nd := xcc.cl, xcc.nd
	for rr := nd[c].down; rr >= xcc.lastItm; rr = nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if nd[nn].color >= 0 {
				uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				nd[dd].up = nn
				nd[uu].down = nn
				nd[cc].itm++
			}
			nn++
		}
	}
	l, r := cl[c].prev, cl[c].next
	cl[r].prev = c
	cl[l].next = cl[r].prev
}

func (xcc *XCC) purify(p int) {
	nd := xcc.nd
	cc := nd[p].itm
	x := nd[p].color
	nd[cc].color = x
	for rr := nd[cc].down; rr >= xcc.lastItm; rr = nd[rr].down {
		if nd[rr].color != x {
			for nn := rr + 1; nn != rr; {
				if nd[nn].color >= 0 {
					uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
					if cc <= 0 {
						nn = uu
						continue
					}
					nd[uu].down = dd
					nd[dd].up = uu
					nd[cc].itm--
				}
				nn++
			}
		} else {
			nd[rr].color = -1
		}
	}
}

func (xcc *XCC) unpurify(p int) {
	nd := xcc.nd
	cc := nd[p].itm
	x := nd[p].color
	for rr := nd[cc].up; rr >= xcc.lastItm; rr = nd[rr].up {
		if nd[rr].color < 0 {
			nd[rr].color = x
		} else {
			for nn := rr - 1; nn != rr; {
				if nd[nn].color >= 0 {
					uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
					if cc <= 0 {
						nn = dd
						continue
					}
					nd[dd].up = nn
					nd[uu].down = nn
					nd[cc].itm++
				}
				nn--
			}
		}
	}
}

func (xcc *XCC) Dance(
	ctx context.Context,
	rd io.Reader,
) (<-chan [][]string, error) {
	if err := inputMatrix(xcc, rd); err != nil {
		return nil, err
	}

	ch := make(chan [][]string)

	go func() {
		defer close(ch)

		var (
			bestItm, curNode   int
			count, level, maxl int
			cl, nd             = xcc.cl, xcc.nd
			choice             = make([]int, maxLevel)
		)

	forward:
		// Set bestItm to the best item for branching: MRV heuristic.
		t := maxNodes
		for k := cl[root].next; t != 0 && k != root; k = cl[k].next {
			if nd[k].itm <= t { // 'itm' is length of node list
				bestItm = k
				t = nd[k].itm
			}
		}

		// Cover bestItm and set choice[level] to nd[bestItm].down
		xcc.cover(bestItm)
		choice[level] = nd[bestItm].down
		curNode = nd[bestItm].down

	advance:
		if curNode == bestItm { // we've tried all options for bestItm
			goto backup
		}

		// Cover all other items of curNode
		for pp := curNode + 1; pp != curNode; {
			cc := nd[pp].itm
			if cc <= 0 {
				pp = nd[pp].up
			} else {
				if nd[pp].color == 0 {
					xcc.cover(cc)
				} else if nd[pp].color > 0 {
					xcc.purify(pp)
				}
				pp++
			}
		}

		if cl[root].next == root {
			if level+1 > maxl {
				if level+1 >= maxLevel {
					log.Fatal(ErrTooManyLevels)
				}
				maxl = level + 1
			}

			count++
			var sol [][]string
			for k := 0; k <= level; k++ {
				sol = append(sol, xcc.getOption(choice[k]))
			}
			select {
			case <-ctx.Done():
				log.Println("Cancelled!")
				return
			case ch <- sol:
			}

			if count >= maxCount {
				return
			}
			goto recover
		}

		level++
		if level > maxl {
			if level >= maxLevel {
				log.Fatal(ErrTooManyLevels)
			}
			maxl = level
		}
		goto forward

	backup:
		xcc.uncover(bestItm)
		if level == 0 {
			return
		}
		level--
		curNode = choice[level]
		bestItm = nd[curNode].itm

	recover:
		for pp := curNode - 1; pp != curNode; {
			cc := nd[pp].itm
			if cc <= 0 {
				pp = nd[pp].down
			} else {
				if nd[pp].color == 0 {
					xcc.uncover(cc)
				} else if nd[pp].color > 0 {
					xcc.unpurify(pp)
				}
				pp--
			}
		}

		choice[level] = nd[curNode].down
		curNode = nd[curNode].down

		goto advance
	}()

	return ch, nil
}
