package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	MysqlDB *sql.DB
)

func InitMysql(host string) error {
	var err error
	MysqlDB, err = sql.Open("mysql", fmt.Sprintf("admin:shht@tcp(%v)/xx?charset=utf8&parseTime=True&loc=Local", host))

	if err != nil {
		return err
	}
	// set pool params
	MysqlDB.SetMaxOpenConns(2000)
	MysqlDB.SetMaxIdleConns(1000)
	MysqlDB.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err = MysqlDB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func UpdateIccid(imei uint64, iccid string) error {
	_, err := MysqlDB.Exec("update device set iccid=? where imei=?", iccid, imei)
	return err
}

func SetVorSwitch(imei uint64, vorSwitch int) error {
	var err error
	bitValue := 1 << 0
	if vorSwitch == 1 {
		_, err = MysqlDB.Exec("update device set switch=switch|? where imei=?", bitValue, imei)
	} else {
		_, err = MysqlDB.Exec("update device set switch=switch&~? where imei=?", bitValue, imei)
	}
	return err
}
