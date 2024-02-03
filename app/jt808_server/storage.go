package main

import (
    "github.com/gocql/gocql"
    "github.com/rayjay214/parser/jt808"
    "strconv"
)

func insertLocation(session *gocql.Session, entity *jt808.T808_0x0200, imei uint64) error {
    date := entity.Time.Format("20060102")
    iDate, _ := strconv.Atoi(date)

    var iLat, iLng uint64
    fLat, _ := entity.Lat.Float64()
    fLng, _ := entity.Lng.Float64()
    iLat = uint64(fLat * float64(1000000))
    iLng = uint64(fLng * float64(1000000))

    query := "insert into t_location(imei, date, time, addr, direction, lat, lng, speed, type) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"

    err := session.Query(query, imei, iDate, entity.Time.Unix(), "", entity.Direction, iLat, iLng, entity.Speed, 0).Exec()
    if err != nil {
        return err
    }

    return nil
}
