package cache

import (
	"sync"
	"time"

	"github.com/ihrk/microbot/internal/unixtime"
)

const minSize = 10

type Cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, v interface{}, dur time.Duration)
}

type node struct {
	v   interface{}
	exp unixtime.Time
}

type cache struct {
	m      sync.RWMutex
	size   int
	values map[string]node
}

func New() Cache {
	return &cache{
		values: make(map[string]node, minSize),
		size:   minSize,
	}
}

func (c *cache) Get(k string) (v interface{}, ok bool) {
	c.m.RLock()

	node, found := c.values[k]
	if found {
		if now := unixtime.Now(); node.exp.After(now) {
			v = node.v
			ok = true
		}
	}

	c.m.RUnlock()

	return
}

func (c *cache) Set(k string, v interface{}, dur time.Duration) {
	c.m.Lock()

	now := unixtime.Now()

	c.cleanup(now)

	c.values[k] = node{v, now.Add(dur)}

	c.m.Unlock()
}

func (c *cache) cleanup(stamp unixtime.Time) {
	if len(c.values) < 2*c.size {
		return
	}

	for key, node := range c.values {
		if node.exp.Before(stamp) {
			delete(c.values, key)
		}
	}

	c.size = max(minSize, len(c.values))
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
