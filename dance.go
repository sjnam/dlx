package dlx

import (
	"log"
	"math/rand"
	"sort"
)

const infty = int(^uint(0) >> 1)

var maxCount = infty

func (x *XCC) visitSolution(level int) [][]string {
	var solution [][]string
	for k := 0; k <= level; k++ {
		var opt []string
		p := x.choice[k]
		for q := p; ; {
			opt = append(opt, x.cl[x.nd[q].itm].name+x.nd[q].scolor)
			q++
			if x.nd[q].itm <= 0 {
				q = x.nd[q].up
			}
			if q == p {
				break
			}
		}
		sort.Strings(opt)
		solution = append(solution, opt)
	}
	return solution
}

func (x *XCC) cover(c int) {
	l, r := x.cl[c].prev, x.cl[c].next
	x.cl[l].next = r
	x.cl[r].prev = l
	for rr := x.nd[c].down; rr >= x.lastItm; rr = x.nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if x.nd[nn].color >= 0 {
				uu := x.nd[nn].up
				dd := x.nd[nn].down
				cc := x.nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				x.nd[uu].down = dd
				x.nd[dd].up = uu
				x.nd[cc].itm--
			}
			nn++
		}
	}
}

func (x *XCC) uncover(c int) {
	for rr := x.nd[c].down; rr >= x.lastItm; rr = x.nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if x.nd[nn].color >= 0 {
				uu := x.nd[nn].up
				dd := x.nd[nn].down
				cc := x.nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				x.nd[dd].up = nn
				x.nd[uu].down = x.nd[dd].up
				x.nd[cc].itm++
			}
			nn++
		}
	}
	l, r := x.cl[c].prev, x.cl[c].next
	x.cl[r].prev = c
	x.cl[l].next = x.cl[r].prev
}

func (x *XCC) purify(p int) {
	cc := x.nd[p].itm
	c := x.nd[p].color
	x.nd[cc].color = c
	for rr := x.nd[cc].down; rr >= x.lastItm; rr = x.nd[rr].down {
		if x.nd[rr].color != c {
			for nn := rr + 1; nn != rr; {
				if x.nd[nn].color >= 0 {
					uu := x.nd[nn].up
					dd := x.nd[nn].down
					cc = x.nd[nn].itm
					if cc <= 0 {
						nn = uu
						continue
					}
					x.nd[uu].down = dd
					x.nd[dd].up = uu
					x.nd[cc].itm--
				}
				nn++
			}
		} else {
			x.nd[rr].color = -1
		}
	}
}

func (x *XCC) unpurify(p int) {
	cc := x.nd[p].itm
	c := x.nd[p].color
	for rr := x.nd[cc].up; rr >= x.lastItm; rr = x.nd[rr].up {
		if x.nd[rr].color < 0 {
			x.nd[rr].color = c
		} else {
			for nn := rr - 1; nn != rr; {
				if x.nd[nn].color >= 0 {
					uu := x.nd[nn].up
					dd := x.nd[nn].down
					cc = x.nd[nn].itm
					if cc <= 0 {
						nn = dd
						continue
					}
					x.nd[dd].up = nn
					x.nd[uu].down = x.nd[dd].up
					x.nd[cc].itm++
				}
				nn--
			}
		}
	}
}

func (x *XCC) Dance() <-chan [][]string {
	ch := make(chan [][]string)

	go func() {
		defer close(ch)

		var p, bestItm, count, curNode, maxl int

		level := 0
	forward:
		t := maxNodes
		for k := x.cl[root].next; t != 0 && k != root; k = x.cl[k].next {
			if x.nd[k].itm <= t {
				if x.nd[k].itm < t {
					bestItm = k
					t = x.nd[k].itm
					p = 1
				} else {
					p++
					if rand.Intn(p) == 0 {
						bestItm = k
					}
				}
			}
		}

		x.cover(bestItm)
		x.choice[level] = x.nd[bestItm].down
		curNode = x.choice[level]

	advance:
		if curNode == bestItm {
			goto backup
		}

		for pp := curNode + 1; pp != curNode; {
			cc := x.nd[pp].itm
			if cc <= 0 {
				pp = x.nd[pp].up
			} else {
				if x.nd[pp].color == 0 {
					x.cover(cc)
				} else if x.nd[pp].color > 0 {
					x.purify(pp)
				}
				pp++
			}
		}

		if x.cl[root].next == root {
			if level+1 > maxl {
				if level+1 >= maxLevel {
					log.Fatal("too many levels")
				}
				maxl = level + 1
			}

			count++
			ch <- x.visitSolution(level)
			if count >= maxCount {
				goto done
			}
			goto recover
		}

		level++
		if level > maxl {
			if level >= maxLevel {
				log.Fatal("too many levels")
			}
			maxl = level
		}
		goto forward

	backup:
		x.uncover(bestItm)
		if level == 0 {
			goto done
		}
		level--
		curNode = x.choice[level]
		bestItm = x.nd[curNode].itm

	recover:
		for pp := curNode - 1; pp != curNode; {
			cc := x.nd[pp].itm
			if cc <= 0 {
				pp = x.nd[pp].down
			} else {
				if x.nd[pp].color == 0 {
					x.uncover(cc)
				} else if x.nd[pp].color > 0 {
					x.unpurify(pp)
				}
				pp--
			}
		}

		x.choice[level] = x.nd[curNode].down
		curNode = x.choice[level]

		goto advance

	done:
	}()

	return ch
}
