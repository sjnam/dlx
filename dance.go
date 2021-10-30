package dlx

import (
	"log"
	"math/rand"
	"sort"
)

func (d *DLX) visitSolution(level int) [][]string {
	var solution [][]string
	for k := 0; k <= level; k++ {
		var opt []string
		p := d.choice[k]
		for q := p; ; {
			opt = append(opt, d.cl[d.nd[q].itm].name+d.nd[q].scolor)
			q++
			if d.nd[q].itm <= 0 {
				q = d.nd[q].up
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

func (d *DLX) cover(c int) {
	l, r := d.cl[c].prev, d.cl[c].next
	d.cl[l].next = r
	d.cl[r].prev = l
	for rr := d.nd[c].down; rr >= d.lastItm; rr = d.nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if d.nd[nn].color >= 0 {
				uu := d.nd[nn].up
				dd := d.nd[nn].down
				cc := d.nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				d.nd[uu].down = dd
				d.nd[dd].up = uu
				d.nd[cc].itm--
			}
			nn++
		}
	}
}

func (d *DLX) uncover(c int) {
	for rr := d.nd[c].down; rr >= d.lastItm; rr = d.nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if d.nd[nn].color >= 0 {
				uu := d.nd[nn].up
				dd := d.nd[nn].down
				cc := d.nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				d.nd[dd].up = nn
				d.nd[uu].down = d.nd[dd].up
				d.nd[cc].itm++
			}
			nn++
		}
	}
	l, r := d.cl[c].prev, d.cl[c].next
	d.cl[r].prev = c
	d.cl[l].next = d.cl[r].prev
}

func (d *DLX) purify(p int) {
	cc := d.nd[p].itm
	x := d.nd[p].color
	d.nd[cc].color = x
	for rr := d.nd[cc].down; rr >= d.lastItm; rr = d.nd[rr].down {
		if d.nd[rr].color != x {
			for nn := rr + 1; nn != rr; {
				if d.nd[nn].color >= 0 {
					uu := d.nd[nn].up
					dd := d.nd[nn].down
					cc = d.nd[nn].itm
					if cc <= 0 {
						nn = uu
						continue
					}
					d.nd[uu].down = dd
					d.nd[dd].up = uu
					d.nd[cc].itm--
				}
				nn++
			}
		} else {
			d.nd[rr].color = -1
		}
	}
}

func (d *DLX) unpurify(p int) {
	cc := d.nd[p].itm
	x := d.nd[p].color
	for rr := d.nd[cc].up; rr >= d.lastItm; rr = d.nd[rr].up {
		if d.nd[rr].color < 0 {
			d.nd[rr].color = x
		} else {
			for nn := rr - 1; nn != rr; {
				if d.nd[nn].color >= 0 {
					uu := d.nd[nn].up
					dd := d.nd[nn].down
					cc = d.nd[nn].itm
					if cc <= 0 {
						nn = dd
						continue
					}
					d.nd[dd].up = nn
					d.nd[uu].down = d.nd[dd].up
					d.nd[cc].itm++
				}
				nn--
			}
		}
	}
}

func (d *DLX) Dance() <-chan [][]string {
	ch := make(chan [][]string)

	go func() {
		defer close(ch)

		var p, bestItm, count, curNode, maxl int

		level := 0
	forward:
		t := maxNodes
		for k := d.cl[root].next; t != 0 && k != root; k = d.cl[k].next {
			if d.nd[k].itm <= t {
				if d.nd[k].itm < t {
					bestItm = k
					t = d.nd[k].itm
					p = 1
				} else {
					p++
					if rand.Intn(p) == 0 {
						bestItm = k
					}
				}
			}
		}

		d.cover(bestItm)
		d.choice[level] = d.nd[bestItm].down
		curNode = d.choice[level]

	advance:
		if curNode == bestItm {
			goto backup
		}

		for pp := curNode + 1; pp != curNode; {
			cc := d.nd[pp].itm
			if cc <= 0 {
				pp = d.nd[pp].up
			} else {
				if d.nd[pp].color == 0 {
					d.cover(cc)
				} else if d.nd[pp].color > 0 {
					d.purify(pp)
				}
				pp++
			}
		}

		if d.cl[root].next == root {
			if level+1 > maxl {
				if level+1 >= maxLevel {
					log.Fatal("too many levels")
				}
				maxl = level + 1
			}

			count++
			ch <- d.visitSolution(level)
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
		d.uncover(bestItm)
		if level == 0 {
			goto done
		}
		level--
		curNode = d.choice[level]
		bestItm = d.nd[curNode].itm

	recover:
		for pp := curNode - 1; pp != curNode; {
			cc := d.nd[pp].itm
			if cc <= 0 {
				pp = d.nd[pp].down
			} else {
				if d.nd[pp].color == 0 {
					d.uncover(cc)
				} else if d.nd[pp].color > 0 {
					d.unpurify(pp)
				}
				pp--
			}
		}

		d.choice[level] = d.nd[curNode].down
		curNode = d.choice[level]

		goto advance

	done:
	}()

	return ch
}
