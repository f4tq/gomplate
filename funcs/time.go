package funcs

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	gotime "time"

	"github.com/hairyhenderson/gomplate/time"
)

var (
	timeNS     *TimeFuncs
	timeNSInit sync.Once
)

// TimeNS -
func TimeNS() *TimeFuncs {
	timeNSInit.Do(func() { timeNS = &TimeFuncs{} })
	return timeNS
}

// AddTimeFuncs -
func AddTimeFuncs(f map[string]interface{}) {
	f["time"] = TimeNS
}

// TimeFuncs -
type TimeFuncs struct{}

// ZoneName - return the local system's time zone's name
func (f *TimeFuncs) ZoneName() string {
	return time.ZoneName()
}

// ZoneOffset - return the local system's time zone offset, in seconds east of UTC
func (f *TimeFuncs) ZoneOffset() int {
	return time.ZoneOffset()
}

// Parse -
func (f *TimeFuncs) Parse(layout, value string) (gotime.Time, error) {
	return gotime.Parse(layout, value)
}

// Now -
func (f *TimeFuncs) Now() gotime.Time {
	return gotime.Now()
}

// Unix - convert UNIX time (in seconds since the UNIX epoch) into a time.Time for further processing
// Takes a string or number (int or float)
func (f *TimeFuncs) Unix(in interface{}) (gotime.Time, error) {
	sec, nsec, err := parseNum(in)
	if err != nil {
		return gotime.Time{}, err
	}
	return gotime.Unix(sec, nsec), nil
}

// convert a number input to a pair of int64s, representing the integer portion and the decimal remainder
// this can handle a string as well as any integer or float type
// precision is at the "nano" level (i.e. 1e+9)
func parseNum(in interface{}) (integral int64, fractional int64, err error) {
	if s, ok := in.(string); ok {
		ss := strings.Split(s, ".")
		if len(ss) > 2 {
			return 0, 0, fmt.Errorf("can not parse '%s' as a number - too many decimal points", s)
		}
		if len(ss) == 1 {
			integral, err := strconv.ParseInt(s, 0, 64)
			return integral, 0, err
		}
		integral, err := strconv.ParseInt(ss[0], 0, 64)
		if err != nil {
			return integral, 0, err
		}
		fractional, err = strconv.ParseInt(ss[1], 0, 64)
		return integral, fractional, err
	}
	if s, ok := in.(fmt.Stringer); ok {
		return parseNum(s.String())
	}
	if i, ok := in.(int); ok {
		return int64(i), 0, nil
	}
	if u, ok := in.(uint64); ok {
		return int64(u), 0, nil
	}
	if f, ok := in.(float64); ok {
		return 0, 0, fmt.Errorf("can not parse floating point number (%f) - use a string instead", f)
	}
	if in == nil {
		return 0, 0, nil
	}
	return 0, 0, nil
}
