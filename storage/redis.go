package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb  *redis.Client
	pipe *redis.Pipeline
)

func InitRedis(host string) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: host,
		DB:   0,
	})
	pipe = Rdb.Pipeline().(*redis.Pipeline)
}

func GetDevice(imei uint64) (map[string]string, error) {
	key := fmt.Sprintf("imei_%v", imei)
	result, err := Rdb.HGetAll(context.Background(), key).Result()
	return result, err
}

func SetRunInfo(imei uint64, info map[string]interface{}) error {
	key := fmt.Sprintf("runinfo_%v", imei)
	state, ok := info["state"]
	if ok && state != nil {
		//状态变化了要记录状态开始时间
		result, err := Rdb.HGetAll(context.Background(), key).Result()
		if err != nil {
			return err
		}

		// 安全地进行类型断言
		if stateStr, ok := state.(string); ok {
			if stateStr != result["state"] {
				info["state_begin_time"] = time.Now()
				UpdateStatus(imei, stateStr)
			}
		}
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()
	return err
}

func GetRunInfo(imei uint64) (map[string]string, error) {
	key := fmt.Sprintf("runinfo_%v", imei)
	result, err := Rdb.HGetAll(context.Background(), key).Result()
	return result, err
}

func SetCmdLog(imei uint64, seqNo uint16, timeid uint64, protocol string) error {
	var key string
	if protocol == "7" || protocol == "9" {
		key = fmt.Sprintf("cmdlog_%v", imei)
	} else {
		key = fmt.Sprintf("cmdlog_%v_%v", imei, seqNo)
	}
	info := map[string]interface{}{
		"timeid": timeid,
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()

	_, err = Rdb.Expire(context.Background(), key, 180*time.Second).Result()

	return err
}

func SetCmdLogZZE(imei uint64, content string, timeid uint64, protocol string) error {
	key := fmt.Sprintf("cmdlog_%v", imei)

	var info map[string]interface{}
	if content == "OFFLINE,1#" {
		info = map[string]interface{}{
			"timeid":       timeid,
			"fake_offline": 1,
		}
	} else if content == "OFFLINE,0#" {
		info = map[string]interface{}{
			"timeid":      timeid,
			"fake_online": 1,
		}
	} else {
		info = map[string]interface{}{
			"timeid": timeid,
		}
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()

	_, err = Rdb.Expire(context.Background(), key, 180*time.Second).Result()

	return err
}

func SetCmdLogC3(imei uint64, content string, timeid uint64, protocol string) error {
	key := fmt.Sprintf("cmdlog_%v", imei)

	var info map[string]interface{}
	if content == "SL GJ" {
		info = map[string]interface{}{
			"timeid":       timeid,
			"fake_offline": 1,
		}
	} else if content == "SL KJ" {
		info = map[string]interface{}{
			"timeid":      timeid,
			"fake_online": 1,
		}
	} else {
		info = map[string]interface{}{
			"timeid": timeid,
		}
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()

	_, err = Rdb.Expire(context.Background(), key, 180*time.Second).Result()

	Rdb.Del(context.Background(), fmt.Sprintf("fakeoff_%v", imei))

	return err
}

func SetCmdLogFakeOnline(imei uint64, content string) error {
	key := fmt.Sprintf("cmdlog_%v", imei)

	info := make(map[string]interface{}, 0)

	if content == "1" {
		info = map[string]interface{}{
			"fake_online": 1,
		}
	} else {
		info = map[string]interface{}{
			"fake_offline": 1,
		}
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()
	_, err = Rdb.Expire(context.Background(), key, 180*time.Second).Result()

	Rdb.Del(context.Background(), fmt.Sprintf("fakeoff_%v", imei))

	return err
}

func SetCmdLogMode(imei uint64, seqNo uint16, timeid uint64, mode string, protocol string) error {
	var key string
	if protocol == "7" || protocol == "8" || protocol == "9" {
		key = fmt.Sprintf("cmdlog_%v", imei)
	} else {
		key = fmt.Sprintf("cmdlog_%v_%v", imei, seqNo)
	}

	info := map[string]interface{}{
		"timeid": timeid,
		"mode":   mode,
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()

	_, err = Rdb.Expire(context.Background(), key, 180*time.Second).Result()

	return err
}

func SetCmdLogShakeValue(imei uint64, seqNo uint16, timeid uint64, value int32) error {
	key := fmt.Sprintf("cmdlog_%v_%v", imei, seqNo)
	info := map[string]interface{}{
		"timeid":      timeid,
		"shake_value": value,
	}

	_, err := Rdb.HSet(context.Background(), key, info).Result()

	_, err = Rdb.Expire(context.Background(), key, 180*time.Second).Result()

	return err
}

func GetCmdLog(imei uint64, seqNo uint16, protocol int) (map[string]string, error) {
	var key string
	if protocol == 7 || protocol == 8 || protocol == 9 {
		key = fmt.Sprintf("cmdlog_%v", imei)
	} else {
		key = fmt.Sprintf("cmdlog_%v_%v", imei, seqNo)
	}

	result, err := Rdb.HGetAll(context.Background(), key).Result()

	return result, err
}

func SetRecordSchedule(imei uint64, schedule float32) error {
	key := fmt.Sprintf("record_schedule_%v", imei)
	_, err := Rdb.Set(context.Background(), key, schedule, 130*time.Second).Result()
	return err
}

func DelRunInfoFields(imei uint64, fields []string) error {
	key := fmt.Sprintf("runinfo_%v", imei)
	_, err := Rdb.HDel(context.Background(), key, fields...).Result()
	return err
}

func SetStartTime(imei uint64) (bool, error) {
	key := fmt.Sprintf("starttime_%v", imei)
	set, err := Rdb.SetNX(context.Background(), key, time.Now(), 0).Result()
	return set, err
}

func CheckFenceSwitch(imei uint64) (bool, error) {
	key := "fenceset"
	isExist, err := Rdb.SIsMember(context.Background(), key, imei).Result()
	return isExist, err
}

func GetFence(imei uint64) (map[string]string, error) {
	key := fmt.Sprintf("fenceinfo_%v", imei)
	result, err := Rdb.HGetAll(context.Background(), key).Result()
	return result, err
}

func SetFence(imei uint64, fenceId string, fenceInfo string) error {
	key := fmt.Sprintf("fenceinfo_%v", imei)
	info := map[string]string{
		fenceId: fenceInfo,
	}
	_, err := Rdb.HSet(context.Background(), key, info).Result()
	return err
}
