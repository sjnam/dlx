package dlx

import "log"

func (d *DLX) visitSolution(ch chan<- Solution, level int, choice []int) {
	var sol Solution
	cl, nd := d.cl, d.nd

	for k := 0; k <= level; k++ {
		var opt []string
		p := choice[k]
		for q := p; ; {
			opt = append(opt, cl[nd[q].itm].name+nd[q].colorName)
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
	ch <- sol
}

func (d *DLX) cover(c int) {
	cl, nd := d.cl, d.nd
	l, r := cl[c].prev, cl[c].next
	cl[l].next = r
	cl[r].prev = l
	for rr := nd[c].down; rr >= d.lastItm; rr = nd[rr].down {
		// hide(rr)
		for nn := rr + 1; nn != rr; {
			if nd[nn].color >= 0 {
				uu := nd[nn].up
				dd := nd[nn].down
				cc := nd[nn].itm
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

func (d *DLX) uncover(c int) {
	cl, nd := d.cl, d.nd
	for rr := nd[c].down; rr >= d.lastItm; rr = nd[rr].down {
		// unhide(rr)
		for nn := rr + 1; nn != rr; {
			if nd[nn].color >= 0 {
				uu := nd[nn].up
				dd := nd[nn].down
				cc := nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				nd[dd].up = nn
				nd[uu].down = nd[dd].up
				nd[cc].itm++
			}
			nn++
		}
	}
	l, r := cl[c].prev, cl[c].next
	cl[r].prev = c
	cl[l].next = cl[r].prev
}

func (d *DLX) purify(p int) {
	nd := d.nd
	cc := nd[p].itm
	x := nd[p].color
	nd[cc].color = x
	for rr := nd[cc].down; rr >= d.lastItm; rr = nd[rr].down {
		if nd[rr].color != x {
			for nn := rr + 1; nn != rr; {
				if nd[nn].color >= 0 {
					uu := nd[nn].up
					dd := nd[nn].down
					cc = nd[nn].itm
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

func (d *DLX) unpurify(p int) {
	nd := d.nd
	cc := nd[p].itm
	x := nd[p].color
	for rr := nd[cc].up; rr >= d.lastItm; rr = nd[rr].up {
		if nd[rr].color < 0 {
			nd[rr].color = x
		} else {
			for nn := rr - 1; nn != rr; {
				if nd[nn].color >= 0 {
					uu := nd[nn].up
					dd := nd[nn].down
					cc = nd[nn].itm
					if cc <= 0 {
						nn = dd
						continue
					}
					nd[dd].up = nn
					nd[uu].down = nd[dd].up
					nd[cc].itm++
				}
				nn--
			}
		}
	}
}

func (d *DLX) Dance() <-chan Solution {
	ch := make(chan Solution)

	go func() {
		defer close(ch)

		var bestItm, curNode, count, maxl int
		cl, nd := d.cl, d.nd
		choice := make([]int, maxLevel)

		level := 0

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
		d.cover(bestItm)
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
					d.cover(cc)
				} else if nd[pp].color > 0 {
					d.purify(pp)
				}
				pp++
			}
		}

		if cl[root].next == root {
			// Visit a solution and goto recover
			if level+1 > maxl {
				if level+1 >= maxLevel {
					log.Fatal(ErrTooManyLevels)
				}
				maxl = level + 1
			}

			count++
			d.visitSolution(ch, level, choice)
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
		d.uncover(bestItm)
		if level == 0 {
			return
		}
		level--
		curNode = choice[level]
		bestItm = nd[curNode].itm

	recover:
		// Uncover all other items of curNode
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
		curNode = nd[curNode].down

		goto advance
	}()

	return ch
}
