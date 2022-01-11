package cooldown

import (
	"sync"
	"time"

	"github.com/ihrk/microbot/internal/unixtime"
)

type Cooldown interface {
	Check() bool
}

type cooldown struct {
	m     sync.Mutex
	d     time.Duration
	exp   unixtime.Time
	gap   int // minimal number of consecutive fails after success
	fails int
}

func New(d time.Duration, gap int) Cooldown {
	return &cooldown{
		exp: unixtime.Now().Add(d),
		d:   d,
		gap: gap,
	}
}

func (cd *cooldown) Check() bool {
	var ok bool

	cd.m.Lock()

	if cd.fails < cd.gap {
		cd.fails++
	} else if now := unixtime.Now(); cd.exp.Before(now) {
		cd.exp = now.Add(cd.d)
		cd.fails = 0
		ok = true
	}

	cd.m.Unlock()

	return ok
}
