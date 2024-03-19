package storage

import "time"

type Location struct {
	Imei      uint64
	Date      int
	Time      int64
	Direction uint16
	Lat       int64
	Lng       int64
	Speed     uint16
	Type      int
	Wgs       string
}

type Record struct {
	Imei     uint64
	Time     int64
	Duration int
	Filename string
	Status   int
}

type Alarm struct {
	Imei      uint64
	Time      time.Time
	Lat       int64
	Lng       int64
	Speed     uint16
	Type      string
	FenceName string
}
