package dlx

import (
	"context"
	"io"
	"log"
)

func (xc *XC) getOption(p int) []string {
	var (
		option []string
		cl, nd = xc.cl, xc.nd
	)
	for q := p; ; {
		option = append(option, cl[nd[q].itm].name)
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

func (xc *XC) cover(c int) {
	cl, nd := xc.cl, xc.nd
	l, r := cl[c].prev, cl[c].next
	cl[l].next = r
	cl[r].prev = l
	for rr := nd[c].down; rr >= xc.lastItm; rr = nd[rr].down {
		for nn := rr + 1; nn != rr; {
			uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
			if cc <= 0 {
				nn = uu
				continue
			}
			nd[uu].down = dd
			nd[dd].up = uu
			nd[cc].itm--
			nn++
		}
	}
}

func (xc *XC) uncover(c int) {
	cl, nd := xc.cl, xc.nd
	for rr := nd[c].down; rr >= xc.lastItm; rr = nd[rr].down {
		for nn := rr + 1; nn != rr; {
			uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
			if cc <= 0 {
				nn = uu
				continue
			}
			nd[dd].up = nn
			nd[uu].down = nn
			nd[cc].itm++
			nn++
		}
	}
	l, r := cl[c].prev, cl[c].next
	cl[r].prev = c
	cl[l].next = cl[r].prev
}

func (xc *XC) Dance(
	ctx context.Context,
	rd io.Reader,
) (<-chan [][]string, error) {
	if err := inputMatrix(xc, rd); err != nil {
		return nil, err
	}

	ch := make(chan [][]string)

	go func() {
		defer close(ch)

		var (
			bestItm, curNode   int
			count, level, maxl int
			cl, nd             = xc.cl, xc.nd
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
		xc.cover(bestItm)
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
				xc.cover(cc)
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
				sol = append(sol, xc.getOption(choice[k]))
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
		xc.uncover(bestItm)
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
				xc.uncover(cc)
				pp--
			}
		}

		choice[level] = nd[curNode].down
		curNode = nd[curNode].down

		goto advance
	}()

	return ch, nil
}
