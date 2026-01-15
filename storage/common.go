package storage

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func SetFakeOnlineState(result map[string]string, imei uint64) {
	//假关机将状态置为离线
	if _, ok := result["fake_offline"]; ok {
		info := map[string]interface{}{
			"comm_time": time.Now(),
			"state":     "1",
		}
		SetRunInfo(imei, info)
		if err != nil {
			log.Warnf("%v update state failed %v", imei, err)
		}
		log.Warnf("%v update fake offline", imei)
		Rdb.HSet(context.Background(), fmt.Sprintf("imei_%v", imei), "fake_online_state", "0")
	}

	//假关机开机
	if _, ok := result["fake_online"]; ok {
		info := map[string]interface{}{
			"comm_time": time.Now(),
			"state":     "3",
		}
		SetRunInfo(imei, info)
		if err != nil {
			log.Warnf("%v update state failed %v", imei, err)
		}
		log.Warnf("%v update fake online", imei)
		Rdb.HSet(context.Background(), fmt.Sprintf("imei_%v", imei), "fake_online_state", "1")
		Rdb.Del(context.Background(), fmt.Sprintf("fakeoff_%v", imei))
	}
}
