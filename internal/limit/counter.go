package limit

import (
	"sync"
	"time"

	"github.com/ihrk/microbot/internal/unixtime"
)

type counter struct {
	m     sync.Mutex
	lim   int
	cur   int
	per   time.Duration
	nodes []node
	off   int
}

type node struct {
	val int
	exp unixtime.Time
}

type Counter interface {
	Add(val int) bool
}

func New(lim int, per time.Duration) Counter {
	return &counter{
		lim: lim,
		per: per,
	}
}

func (l *counter) push(val int, stamp unixtime.Time) {
	if len(l.nodes) == cap(l.nodes) && l.off > 0 {
		n := copy(l.nodes, l.nodes[l.off:])
		l.nodes = l.nodes[:n]

		l.off = 0
	}

	l.cur += val
	l.nodes = append(l.nodes, node{val, stamp})
}

func (l *counter) cleanup(stamp unixtime.Time) {
	for i := l.off; i < len(l.nodes); i++ {
		if l.nodes[i].exp.Before(stamp) {
			l.off++
			l.cur -= l.nodes[i].val
		} else {
			break
		}
	}
}

func (l *counter) Add(val int) (ok bool) {
	l.m.Lock()

	now := unixtime.Now()

	l.cleanup(now)

	if l.cur+val <= l.lim {
		ok = true
		l.push(val, now.Add(l.per))
	}

	l.m.Unlock()

	return
}
