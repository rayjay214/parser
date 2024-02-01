package main

import (
    "github.com/rayjay214/parser/common"
    "time"
)

type ShortRecord struct {
    Imei      uint64
    Writer    common.Writer
    StartTime time.Time
    Schedule  float32
}

type VorRecord struct {
    Imei        uint64
    Writer      common.Writer
    StartTime   time.Time
    EndTime     time.Time
    FirstPacket bool
    PkgCnt      int32
}
