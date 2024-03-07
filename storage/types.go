package storage

type Location struct {
	Imei      uint64
	Date      int
	Time      int64
	Direction uint16
	Lat       uint64
	Lng       uint64
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
