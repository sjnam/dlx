package dlx

import (
	"fmt"
	"log"
)

func (m *MCC) getOption(p, head int) []string {
	var option []string
	cl, nd := m.cl, m.nd
	if (p < m.lastItm && p == head) || (head >= m.lastItm && p == nd[head].itm) {
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

func (m *MCC) cover(c int, deact bool) {
	cl, nd := m.cl, m.nd
	if deact {
		l, r := cl[c].prev, cl[c].next
		cl[l].next, cl[r].prev = r, l
	}
	for rr := nd[c].down; rr >= m.lastItm; rr = nd[rr].down {
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

func (m *MCC) uncover(c int, react bool) {
	cl, nd := m.cl, m.nd
	for rr := nd[c].down; rr >= m.lastItm; rr = nd[rr].down {
		for nn := rr + 1; nn != rr; {
			if nd[nn].color >= 0 {
				uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
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
	if react {
		l, r := cl[c].prev, cl[c].next
		cl[r].prev, cl[l].next = c, c
	}
}

func (m *MCC) purify(p int) {
	nd := m.nd
	cc := nd[p].itm
	x := nd[p].color
	nd[cc].color = x
	for rr := nd[cc].down; rr >= m.lastItm; rr = nd[rr].down {
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
					nd[cc].itm--
				}
				nn++
			}
		} else if rr != p {
			nd[rr].color = -1
		}
	}
}

func (m *MCC) unpurify(p int) {
	nd := m.nd
	cc := nd[p].itm
	x := nd[p].color
	for rr := nd[cc].up; rr >= m.lastItm; rr = nd[rr].up {
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

func (m *MCC) tweak(n, block int) {
	nd := m.nd
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
			nd[cc].itm--
		}
		if nn == n {
			break
		}
		nn++
	}
}

func (m *MCC) untweak(c, x, unblock int) {
	nd := m.nd
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
		m.uncover(c, false)
	}
}

// Dance generates all exact covers will be to repeatedly
// choose an active primary item and to branch on the ways to reduce
// the possibilities for covering that item.
// And we explore all possibilities via depth-first search.
func (m *MCC) Dance() <-chan [][]string {
	ch := make(chan [][]string)

	go func() {
		defer close(ch)

		var bestItm, bestL, bestS, curNode, count, maxl int

		cl, nd := m.cl, m.nd

		choice := make([]int, maxLevel)
		scor := make([]int, maxLevel)
		firstTweak := make([]int, maxLevel)

		level := 0
	forward:
		score := infty
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
			sol := make([][]string, level)
			for k := 0; k < level; k++ {
				pp := choice[k]
				cc := nd[pp].itm
				if pp < m.lastItm {
					cc = pp
				}
				head := firstTweak[k]
				if head == 0 {
					head = nd[cc].down
				}
				sol[k] = m.getOption(pp, head)
			}
			ch <- sol

			count++
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
			m.cover(bestItm, true)
		} else {
			firstTweak[level] = curNode
			if cl[bestItm].bound == 0 {
				m.cover(bestItm, true)
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
			m.tweak(curNode, cl[bestItm].bound)
		} else if cl[bestItm].bound != 0 {
			p, q := cl[bestItm].prev, cl[bestItm].next
			cl[p].next, cl[q].prev = q, p
		}

		if curNode > m.lastItm {
			for pp := curNode + 1; pp != curNode; {
				cc := nd[pp].itm
				if cc <= 0 {
					pp = nd[pp].up
				} else {
					if cc < m.second {
						cl[cc].bound--
						if cl[cc].bound == 0 {
							m.cover(cc, true)
						}
					} else {
						if nd[pp].color == 0 {
							m.cover(cc, true)
						} else if nd[pp].color > 0 {
							m.purify(pp)
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
			m.uncover(bestItm, true)
		} else {
			m.untweak(bestItm, firstTweak[level], cl[bestItm].bound)
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

		if curNode < m.lastItm {
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
				if cc < m.second {
					if cl[cc].bound == 0 {
						m.uncover(cc, true)
					}
					cl[cc].bound++
				} else {
					if nd[pp].color == 0 {
						m.uncover(cc, true)
					} else if nd[pp].color > 0 {
						m.unpurify(pp)
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
