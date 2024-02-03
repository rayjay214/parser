package storage

import "github.com/gocql/gocql"

var (
    cluster *gocql.ClusterConfig
)

func InitCass() {
    cluster = gocql.NewCluster("47.107.69.24") // Replace with source Cassandra node IP
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
