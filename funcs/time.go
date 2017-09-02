package funcs

import (
	"sync"

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
