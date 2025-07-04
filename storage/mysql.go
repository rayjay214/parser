package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var (
	MysqlDB *sql.DB
)

func InitMysql(host string) error {
	var err error
	MysqlDB, err = sql.Open("mysql", fmt.Sprintf("admin:shht@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", host, Conf.Mysql.Dbname))

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

func UpdateMode(imei uint64, mode string) error {
	_, err := MysqlDB.Exec("update device set mode=? where imei=?", mode, imei)
	return err
}

func UpdateStatus(imei uint64, status string) error {
	_, err := MysqlDB.Exec("update device set status=? where imei=?", status, imei)
	return err
}

func UpdateStartTime(imei uint64) error {
	startTime := time.Now().Format("2006-01-02 15:04:05")
	_, err := MysqlDB.Exec("update device set start_time=? where imei=?", startTime, imei)
	return err
}

func UpdateShakeValue(imei uint64, shakeValue int) error {
	_, err := MysqlDB.Exec("update device set shake_value=? where imei=?", shakeValue, imei)
	return err
}

func GetPhoneAndAlarm(imei uint64, alarmType string) (string, string) {
	var bind_phone, dict_value string
	row := MysqlDB.QueryRow("select bind_phone from device where imei=?", imei)
	row.Scan(&bind_phone)

	row = MysqlDB.QueryRow("select dict_label from sys_dict_data where dict_type=? and dict_value=?", "alarm_type", alarmType)
	row.Scan(&dict_value)

	return bind_phone, dict_value
}

func InsertAlarm(alarm Alarm) error {
	log.Infof("insert alarm %v", alarm)
	time := alarm.Time.Format("2006-01-02 15:04:05")
	_, err := MysqlDB.Exec("insert into alarm (imei, time, type, lng, lat, speed, fence_name) values (?,?,?,?,?,?,?)",
		alarm.Imei, time, alarm.Type, alarm.Lng, alarm.Lat, alarm.Speed, alarm.FenceName)
	if err != nil {
		log.Infof("alarm err %v", err)
	}

	//短信推送
	err = CheckAsValue(alarm.Imei, "5")
	if err != nil { //没有增值服务
		log.Infof("%v no push service", alarm.Imei)
		return nil
	}
	bindPhone, alarmName := GetPhoneAndAlarm(alarm.Imei, alarm.Type)
	//推送
	err = Send(bindPhone, fmt.Sprintf("%v", alarm.Imei), alarmName)
	if err == nil {
		UseAsValue(alarm.Imei, "5", 1)
	}

	return err
}

func InsertOfflineAlarm(imei uint64) error {
	alarm := Alarm{
		Imei: imei,
		Time: time.Now(),
		Type: "6",
	}

	runInfo, _ := GetRunInfo(imei)
	fLng, _ := strconv.ParseFloat(runInfo["lng"], 64)
	fLat, _ := strconv.ParseFloat(runInfo["lat"], 64)
	alarm.Lat = int64(fLat * 1000000)
	alarm.Lng = int64(fLng * 1000000)
	return InsertAlarm(alarm)
}

// 查询增值服务可用值
func CheckAsValue(imei uint64, asType string) error {
	var total, used int
	row := MysqlDB.QueryRow("SELECT total, used FROM additional_service WHERE imei=? AND start_time<? AND end_time>? AND deleted_at is NULL "+
		"AND (total-used>0) AND service_type=? order by end_time limit 1",
		imei, time.Now(), time.Now(), asType)

	err := row.Scan(&total, &used)
	log.Infof("%v as type %v total %v used %v", imei, asType, total, used)
	return err
}

// 使用增值服务可用值
func UseAsValue(imei uint64, asType string, usedValue int) error {
	_, err := MysqlDB.Exec("update additional_service set used=used+? WHERE imei=? AND start_time<? AND end_time>? AND deleted_at is NULL "+
		"AND (total-used>0) AND service_type=? order by end_time limit 1",
		usedValue, imei, time.Now(), time.Now(), asType)
	return err
}
