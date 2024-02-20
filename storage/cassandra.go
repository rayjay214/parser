package storage

import (
	"github.com/gocql/gocql"
)

var (
	cluster *gocql.ClusterConfig
)

func InitCass(host string) {
	cluster = gocql.NewCluster(host) // Replace with source Cassandra node IP
	cluster.Keyspace = "gps"
	cluster.Consistency = gocql.LocalOne
}

func GetSession() (*gocql.Session, error) {
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func InsertLocation(loc Location) error {
	query := "insert into t_location(imei, date, time, addr, direction, lat, lng, speed, type) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	session, err := GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Query(query, loc.Imei, loc.Date, loc.Time, "", loc.Direction, loc.Lat, loc.Lng, loc.Speed, 0).Exec()
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
