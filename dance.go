package dlx

import (
	"log"
)

func (d *DLX) solution(level int, choice []int) [][]string {
	var sol [][]string
	cl, nd := d.cl, d.nd
	for k := 0; k <= level; k++ {
		var opt []string
		p := choice[k]
		for q := p; ; {
			opt = append(opt, cl[nd[q].itm].name+nd[q].scolor)
			q++
			if nd[q].itm <= 0 {
				q = nd[q].up
			}
			if q == p {
				break
			}
		}
		sol = append(sol, opt)
	}
	return sol
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
		var choice [maxLevel]int

		cl, nd := d.cl, d.nd
		level := 0
	forward:
		t := maxNodes
		for k := cl[root].next; t != 0 && k != root; k = cl[k].next {
			if nd[k].itm <= t {
				if nd[k].itm < t {
					bestItm = k
					t = nd[k].itm
					p = 1
				} else {
					p++ // this many items achieve the min
				}
			}
		}

		d.cover(bestItm)
		choice[level] = nd[bestItm].down
		curNode = choice[level]

	advance:
		if curNode == bestItm {
			goto backup
		}

		for pp := curNode + 1; pp != curNode; {
			cc := nd[pp].itm
			if cc <= 0 {
				pp = nd[pp].up
			} else {
				if nd[pp].color == 0 {
					d.cover(cc)
				} else if nd[pp].color > 0 {
					d.purify(pp)
				}
				pp++
			}
		}

		if cl[root].next == root {
			if level+1 > maxl {
				if level+1 >= maxLevel {
					log.Fatal("too many levels")
				}
				maxl = level + 1
			}

			count++
			ch <- d.solution(level, choice[:])
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
		curNode = choice[level]
		bestItm = nd[curNode].itm

	recover:
		for pp := curNode - 1; pp != curNode; {
			cc := nd[pp].itm
			if cc <= 0 {
				pp = nd[pp].down
			} else {
				if nd[pp].color == 0 {
					d.uncover(cc)
				} else if nd[pp].color > 0 {
					d.unpurify(pp)
				}
				pp--
			}
		}

		choice[level] = nd[curNode].down
		curNode = choice[level]

		goto advance

	done:
		// do something to finish
		return
	}()

	return ch
}
