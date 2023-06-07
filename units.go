package xo

import "time"

const (
	UnitBytesOfTB = 1000 * UnitBytesOfGB
	UnitBytesOfGB = 1000 * UnitBytesOfMB
	UnitBytesOfMB = 1000 * UnitBytesOfKB
	UnitBytesOfKB = 1000

	UnitBytesOfTiB = 1024 * UnitBytesOfGiB
	UnitBytesOfGiB = 1024 * UnitBytesOfMiB
	UnitBytesOfMiB = 1024 * UnitBytesOfKiB
	UnitBytesOfKiB = 1024
)

var (
	UnitSecondsOfMonth  = 30 * UnitSecondsOfDay
	UnitSecondsOfDay    = 24 * UnitSecondsOfHour
	UnitSecondsOfHour   = 60 * UnitSecondsOfMinute
	UnitSecondsOfMinute = 60 * time.Second
)
