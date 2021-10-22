package queue

import (
	"sync"
	"time"
)

type queue struct {
	sum   int
	off   int
	dur   int64
	m     sync.Mutex
	nodes []node
}

type node struct {
	value int
	stamp int64
}

type Queue interface {
	Push(val int) (sum int)
	Sum() (sum int)
	TryPush(val, lim int) (sum int, ok bool)
}

func NewQueue(d time.Duration) Queue {
	return &queue{
		dur: d.Nanoseconds(),
	}
}

func (q *queue) push(val int, stamp int64) int {
	if len(q.nodes) == cap(q.nodes) && q.off > 0 {
		n := copy(q.nodes, q.nodes[q.off:])
		q.nodes = q.nodes[:n]

		q.off = 0
	}

	q.sum += val
	q.nodes = append(q.nodes, node{val, stamp})

	return q.sum
}

func (q *queue) clean() int64 {
	t := time.Now().UnixNano()

	for i := q.off; i < len(q.nodes); i++ {
		if q.nodes[i].stamp < t {
			q.off++
			q.sum -= q.nodes[i].value
		} else {
			break
		}
	}

	return t + q.dur
}

func (q *queue) Push(val int) (sum int) {
	q.m.Lock()

	stamp := q.clean()

	sum = q.push(val, stamp)

	q.m.Unlock()

	return
}

func (q *queue) Sum() (sum int) {
	q.m.Lock()
	q.clean()
	sum = q.sum
	q.m.Unlock()
	return
}

func (q *queue) TryPush(val, lim int) (sum int, ok bool) {
	q.m.Lock()

	stamp := q.clean()

	if q.sum+val <= lim {
		ok = true
		sum = q.push(val, stamp)
	}

	q.m.Unlock()

	return
}
