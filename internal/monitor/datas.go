package monitor

import "time"

var D int64 = 0
var StartTime = time.Now()

func RegisterData() () {
	D++
}