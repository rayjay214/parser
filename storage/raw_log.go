package storage

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"time"
)

type LogRow struct {
	Imei      uint64
	Direction string //R: 设备->平台, S: 平台->设备
	Message   string
}

var (
	RawLogChannel chan LogRow
)

func InitRawLog() {
	RawLogChannel = make(chan LogRow, 100)
	path := Conf.RawLog.Path
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	log.SetOutput(writer)
	go func() {
		for logInfo := range RawLogChannel {
			log.Printf("#%v#%s#%s", logInfo.Imei, logInfo.Direction, logInfo.Message)
		}
	}()
}
