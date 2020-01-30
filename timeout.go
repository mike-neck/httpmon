package httpmon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeOut struct {
	Amount int
	TimeUnit
}

func (out *TimeOut) AsTime() time.Duration {
	return time.Duration(out.Amount) * out.ToDuration()
}

type TimeUnit int

const (
	NullTimeUnit TimeUnit = iota
	Seconds
	Minutes
)

func (tu TimeUnit) ToDuration() time.Duration {
	switch tu {
	case Seconds:
		return time.Second
	case Minutes:
		return time.Minute
	}
	return time.Duration(0)
}

var timeOutStringPattern = regexp.MustCompile("[smSM]$")

func TimeOutFromString(timeout string) (*TimeOut, error) {
	ts := []byte(timeout)
	locs := timeOutStringPattern.FindIndex(ts)
	if len(locs) != 2 {
		return nil, fmt.Errorf("invalid timeout format: %s", timeout)
	}
	s := locs[0]
	as := ts[:s]
	unit := ts[s:]
	amount, err := strconv.Atoi(string(as))
	if err != nil {
		return nil, fmt.Errorf("invalid number format: %s", string(as))
	}
	if amount <= 0 {
		return nil, fmt.Errorf("invalid negative amount: %d", amount)
	}
	timeUnit, err := timeUnitFromString(string(unit))
	if err != nil {
		return nil, err
	}
	return &TimeOut{
		Amount:   amount,
		TimeUnit: timeUnit,
	}, nil
}

func timeUnitFromString(unit string) (TimeUnit, error) {
	switch strings.ToLower(unit) {
	case "s":
		return Seconds, nil
	case "m":
		return Minutes, nil
	}
	return NullTimeUnit, fmt.Errorf("invalid time unit: %s", unit)
}
