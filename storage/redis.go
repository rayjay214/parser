package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	rdb  *redis.Client
	pipe *redis.Pipeline
)

func InitRedis(host string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
		DB:   0,
	})
	pipe = rdb.Pipeline().(*redis.Pipeline)
}

func GetDevice(imei uint64) (map[string]string, error) {
	key := fmt.Sprintf("imei_%v", imei)
	result, err := rdb.HGetAll(context.Background(), key).Result()
	return result, err
}

func SetRunInfo(imei uint64, info map[string]interface{}) error {
	key := fmt.Sprintf("runinfo_%v", imei)
	if info["state"] != "" {
		//状态变化了要记录状态开始时间
		result, err := rdb.HGetAll(context.Background(), key).Result()
		if err != nil {
			return err
		}
		if info["state"] != result["state"] {
			info["state_begin_time"] = time.Now()
		}
	}

	_, err := rdb.HSet(context.Background(), key, info).Result()
	return err
}

func SetCmdLog(imei uint64, seqNo uint16, timeid uint64) error {
	key := fmt.Sprintf("cmdlog_%v_%v", imei, seqNo)
	info := map[string]interface{}{
		"timeid": timeid,
	}

	_, err := rdb.HSet(context.Background(), key, info).Result()

	_, err = rdb.Expire(context.Background(), key, 180*time.Second).Result()

	return err
}

func GetCmdLog(imei uint64, seqNo uint16) (string, error) {
	key := fmt.Sprintf("cmdlog_%v_%v", imei, seqNo)

	result, err := rdb.HGet(context.Background(), key, "timeid").Result()

	return result, err
}

func SetRecordSchedule(imei uint64, schedule float32) error {
	key := fmt.Sprintf("record_schedule_%v", imei)
	_, err := rdb.Set(context.Background(), key, schedule, 130*time.Second).Result()
	return err
}

func DelRunInfoFields(imei uint64, fields []string) error {
	key := fmt.Sprintf("runinfo_%v", imei)
	_, err := rdb.HDel(context.Background(), key, fields...).Result()
	return err
}
