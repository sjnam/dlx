package dlx

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

func (d *Dancer) getOption(p, head int) []string {
	var (
		option []string
		cl, nd = d.cl, d.nd
	)
	if (p < d.lastItm && p == head) || (head >= d.lastItm && p == nd[head].itm) {
		option = append(option, fmt.Sprintf("null %s", cl[p].name))
	} else {
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
	}
	return option
}

func (d *Dancer) cover(c int, deact bool) {
	cl, nd := d.cl, d.nd
	if deact {
		l, r := cl[c].prev, cl[c].next
		cl[l].next, cl[r].prev = r, l
	}
	d.updates++
	for rr := nd[c].down; rr >= d.lastItm; rr = nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if nd[nn].color >= 0 {
				uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				nd[uu].down = dd
				nd[dd].up = uu
				d.updates++
				nd[cc].itm--
			}
			nn++
		}
	}
}

func (d *Dancer) uncover(c int, react bool) {
	cl, nd := d.cl, d.nd
	for rr := nd[c].down; rr >= d.lastItm; rr = nd[rr].down {
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
	if react {
		l, r := cl[c].prev, cl[c].next
		cl[r].prev, cl[l].next = c, c
	}
}

func (d *Dancer) purify(p int) {
	nd := d.nd
	cc := nd[p].itm
	x := nd[p].color
	nd[cc].color = x
	d.cleansings++
	for rr := nd[cc].down; rr >= d.lastItm; rr = nd[rr].down {
		if nd[rr].color != x {
			for nn := rr + 1; nn != rr; {
				uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				if nd[nn].color >= 0 {
					nd[uu].down = dd
					nd[dd].up = uu
					d.updates++
					nd[cc].itm--
				}
				nn++
			}
		} else if rr != p {
			d.cleansings++
			nd[rr].color = -1
		}
	}
}

func (d *Dancer) unpurify(p int) {
	nd := d.nd
	cc := nd[p].itm
	x := nd[p].color
	for rr := nd[cc].up; rr >= d.lastItm; rr = nd[rr].up {
		if nd[rr].color < 0 {
			nd[rr].color = x
		} else if rr != p {
			for nn := rr - 1; nn != rr; {
				uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
				if cc <= 0 {
					nn = dd
					continue
				}
				if nd[nn].color >= 0 {
					nd[dd].up = nn
					nd[uu].down = nn
					nd[cc].itm++
				}
				nn--
			}
		}
	}
}

func (d *Dancer) tweak(n, block int) {
	nd := d.nd
	nn := n
	if block != 0 {
		nn = n + 1
	}
	for {
		if nd[nn].color >= 0 {
			uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
			if cc <= 0 {
				nn = uu
				continue
			}
			nd[uu].down = dd
			nd[dd].up = uu
			d.updates++
			nd[cc].itm--
		}
		if nn == n {
			break
		}
		nn++
	}
}

func (d *Dancer) untweak(c, x, unblock int) {
	nd := d.nd
	z := nd[c].down
	nd[c].down = x
	rr, qq, k := x, c, 0
	for ; rr != z; qq, rr = rr, nd[rr].down {
		nd[rr].up = qq
		k++
		if unblock != 0 {
			for nn := rr + 1; nn != rr; {
				if nd[nn].color >= 0 {
					uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
					if cc <= 0 {
						nn = uu
						continue
					}
					nd[uu].down = nn
					nd[dd].up = nn
					nd[cc].itm++
				}
				nn++
			}
		}
	}
	nd[rr].up = qq
	nd[c].itm += k
	if unblock == 0 {
		d.uncover(c, false)
	}
}

// Dance generates all exact covers
// Our strategy for generating all exact covers will be to repeatedly
// choose an active primary item and to branch on the ways to reduce
// the possibilities for covering that item.
// And we explore all possibilities via depth-first search.
//
// The neat part of this algorithm is the way the lists are maintained.
// Depth-first search means last-in-first-out maintenance of data structures;
// and it turns out that we need no auxiliary tables to undelete elements from
// lists when backing up. The nodes removed from doubly linked lists remember
// their former neighbors, because we do no garbage collection.
//
// The basic operation is ``covering an item.'' This means removing it
// from the list of items needing to be covered, and ``hiding'' its
// options: removing nodes from other lists whenever they belong to an option of
// a node in this item's list. We cover the chosen item when it has
// |bound=1|.

func (d *Dancer) Dance(
	ctx context.Context,
	rd io.Reader,
) (<-chan [][]string, error) {
	if err := d.inputMatrix(rd); err != nil {
		return nil, err
	}

	ch := make(chan [][]string)

	go func() {
		defer close(ch)

		var (
			bestItm, curNode int
			bestL, bestS     int
			level, maxl      int // maximum level actually reached
			cl, nd           = d.cl, d.nd
			choice           = make([]int, maxLevel)
			scor             = make([]int, maxLevel)
			firstTweak       = make([]int, maxLevel)
		)

	forward:
		d.nodes++
		select {
		case <-ctx.Done():
			log.Println("Cancelled!")
			return
		default:
		}

		score := infty
		for k := cl[root].next; k != root; k = cl[k].next {
			s := cl[k].slack
			if s > cl[k].bound {
				s = cl[k].bound
			}
			t := nd[k].itm + s - cl[k].bound + 1
			p := 0
			if t <= score {
				if t < score || s < bestS || (s == bestS && nd[k].itm > bestL) {
					score = t
					bestItm = k
					bestS = s
					bestL = nd[k].itm
					p = 1
				} else if s == bestS && nd[k].itm == bestL {
					p++
					if rand.Intn(p) == 0 {
						bestItm = k
					}
				}
			}
		}

		if score <= 0 {
			goto backdown
		}
		if score == infty {
			sol := make([][]string, level)
			for k := 0; k < level; k++ {
				pp := choice[k]
				cc := nd[pp].itm
				if pp < d.lastItm {
					cc = pp
				}
				head := firstTweak[k]
				if head == 0 {
					head = nd[cc].down
				}
				sol[k] = d.getOption(pp, head)
			}

			select {
			case <-ctx.Done():
				log.Println("Cancelled!")
				return
			case ch <- sol:
			}

			d.count++
			if d.count >= maxCount {
				return
			}
			goto backdown
		}

		scor[level] = score
		firstTweak[level] = 0

		choice[level] = nd[bestItm].down
		curNode = nd[bestItm].down
		cl[bestItm].bound--

		if cl[bestItm].bound == 0 && cl[bestItm].slack == 0 {
			d.cover(bestItm, true)
		} else {
			firstTweak[level] = curNode
			if cl[bestItm].bound == 0 {
				d.cover(bestItm, true)
			}
		}

	advance:
		if cl[bestItm].bound == 0 && cl[bestItm].slack == 0 {
			if curNode == bestItm {
				goto backup
			}
		} else if nd[bestItm].itm <= cl[bestItm].bound-cl[bestItm].slack {
			goto backup
		} else if curNode != bestItm {
			d.tweak(curNode, cl[bestItm].bound)
		} else if cl[bestItm].bound != 0 {
			p, q := cl[bestItm].prev, cl[bestItm].next
			cl[p].next, cl[q].prev = q, p
		}

		if curNode > d.lastItm {
			for pp := curNode + 1; pp != curNode; {
				cc := nd[pp].itm
				if cc <= 0 {
					pp = nd[pp].up
				} else {
					if cc < d.second {
						cl[cc].bound--
						if cl[cc].bound == 0 {
							d.cover(cc, true)
						}
					} else {
						if nd[pp].color == 0 {
							d.cover(cc, true)
						} else if nd[pp].color > 0 {
							d.purify(pp)
						}
					}
					pp++
				}
			}
		}

		level++
		if level > maxl {
			if level >= maxLevel {
				log.Fatalf("too many levels")
			}
			maxl = level
		}
		goto forward

	backup:
		if cl[bestItm].bound == 0 && cl[bestItm].slack == 0 {
			d.uncover(bestItm, true)
		} else {
			d.untweak(bestItm, firstTweak[level], cl[bestItm].bound)
		}
		cl[bestItm].bound++

	backdown:
		if level == 0 {
			return
		}
		level--
		curNode = choice[level]
		bestItm = nd[curNode].itm
		score = scor[level]

		if curNode < d.lastItm {
			bestItm = curNode
			p, q := cl[bestItm].prev, cl[bestItm].next
			cl[q].prev, cl[p].next = bestItm, bestItm
			goto backup
		}

		for pp := curNode - 1; pp != curNode; {
			cc := nd[pp].itm
			if cc <= 0 {
				pp = nd[pp].down
			} else {
				if cc < d.second {
					if cl[cc].bound == 0 {
						d.uncover(cc, true)
					}
					cl[cc].bound++
				} else {
					if nd[pp].color == 0 {
						d.uncover(cc, true)
					} else if nd[pp].color > 0 {
						d.unpurify(pp)
					}
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

func (d *Dancer) Statistics() {
	s := ""
	if d.count > 1 {
		s = "s"
	}
	fmt.Fprintf(os.Stderr, "Altogether %d solution%s", d.count, s)
	fmt.Fprintf(os.Stderr, " %d updates, %d cleansings, %d nodes.\n",
		d.updates, d.cleansings, d.nodes)
}
