package storage

import (
	log "github.com/sirupsen/logrus"
)

var (
	RawLogChannel chan DeviceLog
)

func InitRawLog() {
	RawLogChannel = make(chan DeviceLog, 10000)
	/*
		path := Conf.RawLog.Path
		writer, _ := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		)
		log.SetOutput(writer)
	*/
	go func() {
		for item := range RawLogChannel {
			//log.Printf("#%v#%s#%s", logInfo.Imei, logInfo.Direction, logInfo.Message)
			err := InsertDeviceLog(item)
			if err != nil {
				log.Warnf("insert device log err %v", err)
			}
		}
	}()
}
