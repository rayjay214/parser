package storage

import (
	log "github.com/sirupsen/logrus"
)

var (
	RawLogChannel chan DeviceLog
)

func InitRawLog() {
	RawLogChannel = make(chan DeviceLog, 10000)
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
