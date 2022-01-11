package unixtime

import (
	"math"
	"time"
)

type Time struct {
	nsec int64
}

func FromTime(tt time.Time) Time {
	return Time{nsec: tt.UnixNano()}
}

func Now() Time {
	return FromTime(time.Now())
}

func Parse(layout, value string) (Time, error) {
	tt, err := time.Parse(layout, value)
	if err != nil {
		return Time{}, err
	}

	return FromTime(tt), nil
}

func (t Time) Add(d time.Duration) Time {
	t.nsec += int64(d)
	return t
}

func (t Time) Sub(u Time) time.Duration {
	d := time.Duration(t.nsec - u.nsec)
	if d < 0 && t.nsec > u.nsec {
		return math.MaxInt64
	}
	if d > 0 && t.nsec < u.nsec {
		return math.MinInt64
	}
	return d
}

func (t Time) Time() time.Time {
	return time.Unix(0, t.nsec)
}

func (t Time) Before(u Time) bool {
	return t.nsec < u.nsec
}

func (t Time) After(u Time) bool {
	return t.nsec > u.nsec
}

func (t Time) Equal(u Time) bool {
	return t.nsec == u.nsec
}

func (t Time) IsZero() bool {
	return t.nsec == 0
}

func (t Time) String() string {
	return t.Time().String()
}

func (t Time) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var err error
	*t, err = Parse(`"`+time.RFC3339+`"`, string(data))

	return err
}
