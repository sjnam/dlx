package dlx

import (
	"fmt"
	"io"
	"os"
	"time"
)

func (m *MCC) option(p, head, score int) Option {
	var opt Option
	cl, nd := m.cl, m.nd

	if (p < m.lastItm && p == head) || (head >= m.lastItm && p == nd[head].itm) {
		opt = append(opt, fmt.Sprintf("null %s", cl[p].name))
	} else {
		q := p + 1
		for q != p {
			if nd[q].itm <= 0 {
				q = nd[q].up
				break
			}
			q++
		}
		for nd[q].itm > 0 {
			opt = append(opt,
				fmt.Sprintf("%s%s", cl[nd[q].itm].name, nd[q].colorName))
			q++
		}
	}

	if m.Debug {
		fmt.Fprintf(os.Stderr, "%s", opt)
		k := 1
		for q := head; q != p; k++ {
			if p >= m.lastItm && q == nd[p].itm {
				fmt.Fprintln(os.Stderr, "(?)")
				return opt
			} else {
				q = nd[q].down
			}
		}
		fmt.Fprintf(os.Stderr, " (%d of %d)\n", k, score)
	}

	return opt
}

func (m *MCC) hide(rr int) {
	nd := m.nd
	for nn := rr + 1; nn != rr; {
		if nd[nn].color >= 0 {
			uu, dd, cc := nd[nn].up, nd[nn].down, nd[nn].itm
			if cc <= 0 {
				nn = uu
				continue
			}
			nd[uu].down = dd
			nd[dd].up = uu
			m.updates++
			nd[cc].itm--
		}
		nn++
	}
}

func (m *MCC) cover(c int, deact bool) {
	cl, nd := m.cl, m.nd
	if deact {
		l, r := cl[c].prev, cl[c].next
		cl[l].next, cl[r].prev = r, l
	}
	m.updates++
	for rr := nd[c].down; rr >= m.lastItm; rr = nd[rr].down {
		m.hide(rr)
	}
}

func (m *MCC) unhide(rr int) {
	nd := m.nd
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

func (m *MCC) uncover(c int, react bool) {
	cl, nd := m.cl, m.nd
	for rr := nd[c].down; rr >= m.lastItm; rr = nd[rr].down {
		m.unhide(rr)
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
	m.cleansings++
	for rr := nd[cc].down; rr >= m.lastItm; rr = nd[rr].down {
		if nd[rr].color != x {
			m.hide(rr)
		} else if rr != p {
			m.cleansings++
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
			m.unhide(rr)
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
			m.updates++
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
			m.unhide(rr)
		}
	}
	nd[rr].up = qq
	nd[c].itm += k
	if unblock == 0 {
		m.uncover(c, false)
	}
}

// Dance generates all exact covers
func (m *MCC) Dance(rd io.Reader) Result {
	if err := m.inputMatrix(rd); err != nil {
		panic(err)
	}

	heartbeat := make(chan string)
	solStream := make(chan []Option)

	go func() {
		defer close(heartbeat)
		defer close(solStream)

		var (
			bestItm, curNode int
			bestL, bestS     int
			level, maxl      int // maximum level actually reached
			cl, nd           = m.cl, m.nd
			choice, scor     [maxLevel]int
			firstTweak       [maxLevel]int
			count, nodes     int
		)

		pulse := time.NewTicker(m.PulseInterval)
		sendPulse := func() {
			select {
			case heartbeat <- fmt.Sprintf("L(%d/%d): %d sols so far",
				level, maxl, count):
			default:
			}
		}

	forward:
		nodes++
		select {
		case <-m.ctx.Done():
			return
		case <-pulse.C:
			sendPulse()
		default:
		}
		// Set bestItm to the best item for branching,
		// and let score be its branching degree
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

		if score <= 0 { // not enough options left in this item
			goto backdown
		}
		if score == infty { // Visit a solution and goto backdown
			sol := make([]Option, level)
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
				sol[k] = m.option(pp, head, scor[k])
			}

			select {
			case <-m.ctx.Done():
				goto done
			case <-pulse.C:
				sendPulse()
			case solStream <- sol:
			}

			count++
			if count >= maxCount {
				goto done
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

	advance: // If curNode is off limits, goto backup; also tweak if needed
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

		if m.Debug {
			fmt.Fprintf(os.Stderr, "L%d: ", level)
			if cl[bestItm].bound == 0 && cl[bestItm].slack == 0 {
				m.option(curNode, nd[bestItm].down, score)
			} else {
				m.option(curNode, firstTweak[level], score)
			}
		}

		if curNode > m.lastItm {
			// Cover or partially cover all other items of curNode's option
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

		// Increase level and goto forward
		level++
		if level > maxl {
			if level >= maxLevel {
				panic("too many levels")
			}
			maxl = level
		}
		goto forward

	backup: // Restore the original state of bestItm
		if cl[bestItm].bound == 0 && cl[bestItm].slack == 0 {
			m.uncover(bestItm, true)
		} else {
			m.untweak(bestItm, firstTweak[level], cl[bestItm].bound)
		}
		cl[bestItm].bound++

	backdown:
		if level == 0 {
			goto done
		}
		level--
		curNode = choice[level]
		bestItm = nd[curNode].itm
		score = scor[level]

		if curNode < m.lastItm {
			// Reactivate bestItm and goto backup
			bestItm = curNode
			p, q := cl[bestItm].prev, cl[bestItm].next
			cl[q].prev, cl[p].next = bestItm, bestItm // reactivate bestItm
			goto backup
		}

		// Uncover or partially uncover all other items of curNode's option
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

	done:
		if m.Debug {
			s := ""
			if count > 1 {
				s = "s"
			}
			fmt.Fprintf(os.Stderr, "Altogether %d solution%s", count, s)
			fmt.Fprintf(os.Stderr, " %d updates, %d cleansings, %d nodes.\n",
				m.updates, m.cleansings, nodes)
		}
	}()

	return Result{
		Solutions: solStream,
		Heartbeat: heartbeat,
	}
}
