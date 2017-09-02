package time

import (
	"time"
)

// ZoneName -
func ZoneName() string {
	n, _ := time.Now().Zone()
	return n
}

// ZoneOffset -
func ZoneOffset() int {
	_, o := time.Now().Zone()
	return o
}
