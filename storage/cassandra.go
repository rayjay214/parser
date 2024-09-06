package storage

import (
	"github.com/gocql/gocql"
	"time"
)

var (
	cluster *gocql.ClusterConfig
)

func InitCass(host string) {
	cluster = gocql.NewCluster(host) // Replace with source Cassandra node IP
	cluster.Keyspace = "gps"
	cluster.Consistency = gocql.LocalOne
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 2}
}

func GetSession() (*gocql.Session, error) {
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func InsertLocation(loc Location) error {
	query := "insert into t_location(imei, date, time, direction, lat, lng, speed, type) values (?, ?, ?, ?, ?, ?, ?, ?)"
	session, err := GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Query(query, loc.Imei, loc.Date, loc.Time, loc.Direction, loc.Lat, loc.Lng, loc.Speed, loc.Type).Exec()
	if err != nil {
		return err
	}

	return nil
}

func InsertRecord(record Record) error {
	query := "insert into t_record(imei, time, duration, filename, status) values (?, ?, ?, ?, ?)"
	session, err := GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Query(query, record.Imei, record.Time, record.Duration, record.Filename, record.Status).Exec()
	if err != nil {
		return err
	}

	return nil
}

func UpdateCmdResponse(imei uint64, timeid uint64, response string) error {
	query := "update t_cmd set response=? where imei=? and time=?"
	session, err := GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Query(query, response, imei, timeid).Exec()
	if err != nil {
		return err
	}

	return nil
}

func InsertDeviceLog(log DeviceLog) error {
	query := "insert into t_device_log(imei, time, raw, type) values (?, ?, ?, ?)"
	session, err := GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Query(query, log.Imei, log.Time, log.Raw, log.Type).Exec()
	if err != nil {
		return err
	}

	return nil
}
