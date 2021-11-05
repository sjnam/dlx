package dlx

import "log"

func (d *DLX) visitSolution(ch chan<- Solution,
	level int, choice, firstTweak []int) {
	var sol Solution
	cl, nd := d.cl, d.nd

	for k := 0; k < level; k++ {
		var opt []string
		pp := choice[k]
		cc := nd[pp].itm
		if pp < d.lastItm {
			cc = pp
		}

		head := firstTweak[k]
		if head == 0 {
			head = nd[cc].down
		}

		if (pp < d.lastItm && pp == head) ||
			(head >= d.lastItm && pp == nd[head].itm) {
			opt = append(opt, "null "+cl[pp].name)
		} else {
			for q := pp; ; {
				opt = append(opt, cl[nd[q].itm].name+nd[q].colorName)
				q++
				if nd[q].itm <= 0 {
					q = nd[q].up
				}
				if q == pp {
					break
				}
			}
		}
		sol = append(sol, opt)
	}
	ch <- sol
}

func (d *DLX) cover(c, deact int) {
	cl, nd := d.cl, d.nd
	if deact != 0 {
		l, r := cl[c].prev, cl[c].next
		cl[l].next = r
		cl[r].prev = l
	}
	for rr := nd[c].down; rr >= d.lastItm; rr = nd[rr].down {
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

func (d *DLX) uncover(c, react int) {
	cl, nd := d.cl, d.nd
	for rr := nd[c].down; rr >= d.lastItm; rr = nd[rr].down {
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
	if react != 0 {
		l, r := cl[c].prev, cl[c].next
		cl[r].prev = c
		cl[l].next = c
	}
}

func (d *DLX) purify(p int) {
	nd := d.nd
	cc := nd[p].itm
	x := nd[p].color
	nd[cc].color = x
	for rr := nd[cc].down; rr >= d.lastItm; rr = nd[rr].down {
		if nd[rr].color != x {
			for nn := rr + 1; nn != rr; {
				uu := nd[nn].up
				dd := nd[nn].down
				cc = nd[nn].itm
				if cc <= 0 {
					nn = uu
					continue
				}
				if nd[nn].color >= 0 {
					nd[uu].down = dd
					nd[dd].up = uu
					nd[cc].itm--
				}
				nn++
			}
		} else if rr != p {
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
		} else if rr != p {
			for nn := rr - 1; nn != rr; {
				uu := nd[nn].up
				dd := nd[nn].down
				cc = nd[nn].itm
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

func (d *DLX) tweak(n, block int) {
	nd := d.nd
	nn := n
	if block != 0 {
		nn = n + 1
	}
	for {
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
		if nn == n {
			break
		}
		nn++
	}
}

func (d *DLX) untweak(c, x, unblock int) {
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
					uu := nd[nn].up
					dd := nd[nn].down
					cc := nd[nn].itm
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
		d.uncover(c, 0)
	}
}

func (d *DLX) Dance() <-chan Solution {
	ch := make(chan Solution)

	go func() {
		defer close(ch)

		var bestItm, bestL, bestS, curNode, count, maxl, score int

		cl, nd := d.cl, d.nd

		choice := make([]int, maxLevel)
		scor := make([]int, maxLevel)
		firstTweak := make([]int, maxLevel)

		level := 0
	forward:
		score = infty
		for k := cl[root].next; k != root; k = cl[k].next {
			s := cl[k].slack
			if s > cl[k].bound {
				s = cl[k].bound
			}
			t := nd[k].itm + s - cl[k].bound + 1
			if t <= score {
				if t < score || s < bestS || (s == bestS && nd[k].itm > bestL) {
					score = t
					bestItm = k
					bestS = s
					bestL = nd[k].itm
				}
			}
		}

		if score <= 0 {
			goto backdown
		}
		if score == infty {
			count++
			d.visitSolution(ch, level, choice, firstTweak)
			if count >= maxCount {
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
			d.cover(bestItm, 1)
		} else {
			firstTweak[level] = curNode
			if cl[bestItm].bound == 0 {
				d.cover(bestItm, 1)
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
			cl[p].next = q
			cl[q].prev = p
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
							d.cover(cc, 1)
						}
					} else {
						if nd[pp].color == 0 {
							d.cover(cc, 1)
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
				log.Fatal(ErrTooManyLevels)
			}
			maxl = level
		}
		goto forward

	backup:
		if cl[bestItm].bound == 0 && cl[bestItm].slack == 0 {
			d.uncover(bestItm, 1)
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
			cl[q].prev = bestItm
			cl[p].next = bestItm
			goto backup
		}

		for pp := curNode - 1; pp != curNode; {
			cc := nd[pp].itm
			if cc <= 0 {
				pp = nd[pp].down
			} else {
				if cc < d.second {
					if cl[cc].bound == 0 {
						d.uncover(cc, 1)
					}
					cl[cc].bound++
				} else {
					if nd[pp].color == 0 {
						d.uncover(cc, 1)
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

	return ch
}
