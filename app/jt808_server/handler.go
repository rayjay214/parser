package main

import (
	"bytes"
	"context"
	"fmt"
	geo "github.com/kellydunn/golang-geo"
	"github.com/qichengzx/coordtransform"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808"
	"github.com/rayjay214/parser/protocol/jt808/extra"
	"github.com/rayjay214/parser/server_base/jt808_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"time"
)

func handle0100(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0100)
	log.Infof("handle 0100 %v", entity)

	session.ReplyRegister(message)
}

func handle0102(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0102)
	log.Infof("handle 0102 %v", entity)

	deviceInfo, _ := storage.GetDevice(message.Header.Imei)
	if len(deviceInfo) == 0 {
		log.Warnf("imei %v not exist", message.Header.Imei)
		session.Close()
		return
	}
	session.Protocol, _ = strconv.Atoi(deviceInfo["protocol"])

	//假关机状态下不更新状态
	if session.Protocol == 7 {
		fakeOnline, _ := storage.Rdb.HGet(context.Background(), fmt.Sprintf("imei_%v", session.ID()), "fake_online").Result()
		if fakeOnline == "0" {
			session.Reply(message, jt808.T808_0x8100_ResultSuccess)
			return
		}
	}

	//C3假关机状态下不更新状态
	if session.Protocol == 8 {
		fakeOnline, _ := storage.Rdb.HGet(context.Background(), fmt.Sprintf("fakeoff_%v", session.ID()), "fake_heartbeat").Result()
		if fakeOnline == "1" {
			session.Reply(message, jt808.T808_0x8100_ResultSuccess)
			return
		}
	}

	set, err := storage.SetStartTime(session.ID())
	if err == nil && set {
		storage.UpdateStartTime(session.ID())
	}

	info := map[string]interface{}{
		//"status":    "2", //maybe useless
		"state":     "2",
		"comm_time": time.Now(),
	}
	storage.SetRunInfo(message.Header.Imei, info)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0002(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0002)

	//假关机状态下不更新状态
	if session.Protocol == 7 {
		lastRunInfo, _ := storage.GetRunInfo(session.ID())
		if lastRunInfo["state"] == "3" {
			session.Reply(message, jt808.T808_0x8100_ResultSuccess)
			return
		}
	}

	info := map[string]interface{}{
		"comm_time": time.Now(),
		"state":     "3",
	}

	for _, ext := range entity.Extras {
		switch ext.ID() {
		case extra.Extra_0x04{}.ID():
			v := ext.(*extra.Extra_0x04).Value().(extra.Extra_0x04_Value)
			info["power"] = v.Power
			info["acc_power"] = v.Status
		case extra.Extra_0x30{}.ID():
			info["signal"] = ext.(*extra.Extra_0x30).Value()
		case extra.Extra_0x31{}.ID():
			info["satellite"] = ext.(*extra.Extra_0x31).Value()
		case extra.Extra_0xf0{}.ID():
			log.Infof("voltage is %v", ext.(*extra.Extra_0xf0).Value())
			info["voltage"] = fmt.Sprintf("%.2fV", float32(ext.(*extra.Extra_0xf0).Value().(uint16)/100))
		case extra.Extra_0xf5{}.ID():
			info["gprs_type"] = ext.(*extra.Extra_0xf5).Value()
		case extra.Extra_0xe5{}.ID():
			v := ext.(*extra.Extra_0xe5).Value().(byte)
			if v == 1 {
				info["state"] = "2"
			} else {
				info["state"] = "3"
			}
		}
	}

	storage.SetRunInfo(message.Header.Imei, info)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0f02(session *jt808_base.Session, message *jt808.Message) {
	key := fmt.Sprintf("fakeoff_%v", message.Header.Imei)
	info := map[string]interface{}{
		"fake_heartbeat": 1,
	}
	storage.Rdb.HSet(context.Background(), key, info).Result()
	storage.Rdb.Expire(context.Background(), key, 600*time.Second).Result()

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0808(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0808)
	//session.Protocol = 2 //2013
	log.Infof("handle 0808 %v", entity)
}

// 处理上报位置
func handle0200(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0200)
	log.Infof("%v handle 0200", session.ID())
	handleLocation(message.Header.Imei, entity, session.Protocol)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0704(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0704)
	log.Infof("%v handle 0704", session.ID())

	for _, item := range entity.Items {
		handleLocation(message.Header.Imei, &item, session.Protocol)
	}

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0201(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0201)
	log.Infof("handle 0201 %v", entity)

	handleLocation(message.Header.Imei, &entity.Result, session.Protocol)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handleLocation(imei uint64, entity *jt808.T808_0x0200, protocol int) {
	date := entity.Time.Format("20060102")
	iDate, _ := strconv.Atoi(date)

	info := map[string]interface{}{
		"comm_time": time.Now(),
		"direction": entity.Direction,
	}
	locTypeBase := 0

	if protocol == 7 {
		if entity.Speed > 0 {
			info["state"] = "2"
		} else {
			info["state"] = "3"
		}
	}

	for _, ext := range entity.Extras {
		switch ext.ID() {
		case extra.Extra_0x01{}.ID():
			info["distance"] = ext.(*extra.Extra_0x01).Value()
		case extra.Extra_0x04{}.ID():
			v := ext.(*extra.Extra_0x04).Value().(extra.Extra_0x04_Value)
			if protocol != 7 {
				info["power"] = v.Power
			}
			info["acc_power"] = v.Status
		case extra.Extra_0xe4{}.ID():
			v := ext.(*extra.Extra_0xe4).Value().(extra.Extra_0xe4_Value)
			info["power"] = v.Power
			info["acc_power"] = v.Status
		case extra.Extra_0x05{}.ID():
			if protocol == 1 {
				v := ext.(*extra.Extra_0x05).Value().(byte)
				if v == 1 {
					info["state"] = "2"
				} else {
					info["state"] = "3"
					locTypeBase = 3
				}
			}
		case extra.Extra_0x30{}.ID():
			info["signal"] = ext.(*extra.Extra_0x30).Value()
		case extra.Extra_0x31{}.ID():
			info["satellite"] = ext.(*extra.Extra_0x31).Value()
		case extra.Extra_0xe7{}.ID():
			v := ext.(*extra.Extra_0xe7).Value().(extra.Extra_0xe7_Value)
			//log.Infof("e7 status is %v", v)
			if protocol == 2 {
				if v.SleepCheckWay == 1 {
					if v.SleepStatus == 1 {
						info["state"] = "3"
						locTypeBase = 3
					} else {
						info["state"] = "2"
					}
				} else {
					if entity.Status.GetAccState() {
						info["state"] = "2"
					} else {
						info["state"] = "3"
						locTypeBase = 3
					}
				}
			}
		case extra.Extra_0xf0{}.ID():
			log.Infof("voltage is %v", ext.(*extra.Extra_0xf0).Value())
			info["voltage"] = fmt.Sprintf("%.2fV", float32(ext.(*extra.Extra_0xf0).Value().(uint16)/100))
		case extra.Extra_0xf5{}.ID():
			info["gprs_type"] = ext.(*extra.Extra_0xf5).Value()
		case extra.Extra_0xe5{}.ID():
			v := ext.(*extra.Extra_0xe5).Value().(byte)
			if v == 1 {
				info["state"] = "2"
			} else {
				info["state"] = "3"
				locTypeBase = 3
			}
		case extra.Extra_0x2b{}.ID():
			v := ext.(*extra.Extra_0x2b).Value().(extra.Extra_0x2b_Value)
			power := CalculateBatteryPercent(int(v.Voltage1 * 10))
			info["power"] = fmt.Sprintf("%v", power)
		}
	}

	lastRunInfo, _ := storage.GetRunInfo(imei)

	var loc storage.Location
	if entity.Status.Positioning() {
		var iLat, iLng int64
		fLat, _ := entity.Lat.Float64()
		fLng, _ := entity.Lng.Float64()
		iLat = int64(fLat * float64(1000000))
		iLng = int64(fLng * float64(1000000))

		loc = storage.Location{
			Imei:      imei,
			Date:      iDate,
			Time:      entity.Time.Unix(),
			Direction: entity.Direction,
			Lat:       iLat,
			Lng:       iLng,
			Speed:     entity.Speed,
			Type:      0 + locTypeBase,
			Wgs:       "",
		}

		err := storage.InsertLocation(loc)
		if err != nil {
			log.Warnf("insert location err %v", err)
		}
		info["lat"] = fLat
		info["lng"] = fLng
		info["loc_type"] = 0 + locTypeBase
		info["loc_time"] = entity.Time
		info["gps_lat"] = fLat
		info["gps_lng"] = fLng

		if lastRunInfo["gps_lat"] != "" {
			lastLat, _ := strconv.ParseFloat(lastRunInfo["gps_lat"], 64)
			lastLng, _ := strconv.ParseFloat(lastRunInfo["gps_lng"], 64)
			lastDayDistance, _ := strconv.Atoi(lastRunInfo["day_distance"])
			lastTotalDistance, _ := strconv.Atoi(lastRunInfo["total_distance"])
			point1 := geo.NewPoint(lastLat, lastLng)
			point2 := geo.NewPoint(fLat, fLng)
			distance := int(point1.GreatCircleDistance(point2) * 1000) //m
			info["day_distance"] = lastDayDistance + distance
			info["total_distance"] = lastTotalDistance + distance
		}
	} else {
		var lbsResp LbsResp
		err := getLbsLocation(entity, &lbsResp, protocol)
		if err != nil {
			return
		}
		if lbsResp.Lng == 0 {
			log.WithFields(log.Fields{
				"imei": imei,
			}).Info("get lbs location failed")
			return
		}

		wgsLng, wgsLat := coordtransform.GCJ02toWGS84(lbsResp.Lng, lbsResp.Lat)

		loc = storage.Location{
			Imei:      imei,
			Date:      iDate,
			Time:      entity.Time.Unix(),
			Direction: entity.Direction,
			Lat:       int64(wgsLat * 1000000),
			Lng:       int64(wgsLng * 1000000),
			Speed:     entity.Speed,
			Type:      lbsResp.LocType + locTypeBase,
			Wgs:       "",
		}

		//wifi上报更新最后位置（基站的不更新），但是如果是首次上报位置，还是要更新基站
		//if lbsResp.LocType == 2 || lastRunInfo["lat"] == "" {
		info["lat"] = wgsLat
		info["lng"] = wgsLng
		info["loc_type"] = lbsResp.LocType + locTypeBase
		//}

		info["loc_time"] = entity.Time

		err = storage.InsertLocation(loc)
		if err != nil {
			log.Warnf("insert location err %v", err)
		}
	}
	storage.SetRunInfo(imei, info)

	checkAlarm(entity, &loc, protocol)

	checkFence(&loc)
}

func checkAlarm(entity *jt808.T808_0x0200, loc *storage.Location, protocol int) {
	fGetAlarmbit := func(value uint32, offset int) uint32 {
		return uint32(value) & (1 << offset) >> offset
	}

	if entity.Alarm == 0 {
		return
	}
	alarm := storage.Alarm{
		Imei:  loc.Imei,
		Time:  time.Unix(loc.Time, 0),
		Lat:   loc.Lat,
		Lng:   loc.Lng,
		Speed: loc.Speed,
	}

	if fGetAlarmbit(entity.Alarm, 0) == 1 { //sos
		alarm.Type = "4"
		storage.InsertAlarm(alarm)
	}
	if fGetAlarmbit(entity.Alarm, 1) == 1 { //超速
		alarm.Type = "5"
		storage.InsertAlarm(alarm)
	}
	if fGetAlarmbit(entity.Alarm, 7) == 1 { //低电
		alarm.Type = "3"
		storage.InsertAlarm(alarm)
	}
	if fGetAlarmbit(entity.Alarm, 9) == 1 { //震动
		alarm.Type = "1"
		storage.InsertAlarm(alarm)
	}

}

func handle1007(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x1007)
	log.Infof("handle 1007 %v", entity)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1107(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x1107)
	log.Infof("handle 1107 %v", entity)

	deviceInfo, err := storage.GetDevice(session.ID())
	if err != nil {
		session.Reply(message, jt808.T808_0x8100_ResultSuccess)
		return
	}
	if iccid, ok := deviceInfo["iccid"]; ok {
		if iccid != entity.Iccid {
			log.Infof("%v update iccid from %v to %v", session.ID(), iccid, entity.Iccid)
			err = storage.UpdateIccid(session.ID(), entity.Iccid)
			if err != nil {
				log.Warnf("%v update iccid failed %v", session.ID(), err)
			}
		}
	}

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1300(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x1300)
	result, err := storage.GetCmdLog(session.ID(), entity.AckSeqNo, session.Protocol)
	if err != nil {
		return
	}

	log.Infof("%v handle 1300, content:%v, result:%v", message.Header.Imei, entity.Content, result)
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}

	//同步定位模式
	if mode, ok := result["mode"]; ok {
		err = storage.UpdateMode(session.ID(), mode)
		if err != nil {
			log.Warnf("%v update mode failed %v", session.ID(), err)
		}
	}

	//假关机将状态置为离线
	if _, ok := result["fake_offline"]; ok {
		info := map[string]interface{}{
			"comm_time": time.Now(),
			"state":     "1",
		}
		storage.SetRunInfo(message.Header.Imei, info)
		if err != nil {
			log.Warnf("%v update state failed %v", session.ID(), err)
		}
		log.Warnf("%v update fake offline", session.ID())
		storage.Rdb.HSet(context.Background(), fmt.Sprintf("imei_%v", session.ID()), "fake_online", "0")
	}

	//假关机开机
	if _, ok := result["fake_online"]; ok {
		info := map[string]interface{}{
			"comm_time": time.Now(),
			"state":     "3",
		}
		storage.SetRunInfo(message.Header.Imei, info)
		if err != nil {
			log.Warnf("%v update state failed %v", session.ID(), err)
		}
		log.Warnf("%v update fake online", session.ID())
		storage.Rdb.HSet(context.Background(), fmt.Sprintf("imei_%v", session.ID()), "fake_online", "1")
		storage.Rdb.Del(context.Background(), fmt.Sprintf("fakeoff_%v", session.ID()))
	}

	err = storage.UpdateCmdResponse(session.ID(), timeid, entity.Content)
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handle6006(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x6006)
	log.Infof("%v handle 6006 %v", session.ID(), entity)
	result, err := storage.GetCmdLog(session.ID(), entity.AckSeqNo, session.Protocol)
	if err != nil {
		return
	}
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}

	//同步定位模式
	if mode, ok := result["mode"]; ok {
		err = storage.UpdateMode(session.ID(), mode)
		if err != nil {
			log.Warnf("%v update mode failed %v", session.ID(), err)
		}
	}

	//同步震动报警阈值
	if value, ok := result["shake_value"]; ok {
		shakeValue, _ := strconv.Atoi(value)
		err = storage.UpdateShakeValue(session.ID(), shakeValue)
		if err != nil {
			log.Warnf("%v update mode failed %v", session.ID(), err)
		}
	}

	//假关机将状态置为离线
	if _, ok := result["fake_offline"]; ok {
		info := map[string]interface{}{
			"comm_time": time.Now(),
			"state":     "1",
		}
		storage.SetRunInfo(message.Header.Imei, info)
		if err != nil {
			log.Warnf("%v update state failed %v", session.ID(), err)
		}
		log.Warnf("%v update fake offline", session.ID())
		storage.Rdb.HSet(context.Background(), fmt.Sprintf("imei_%v", session.ID()), "fake_online", "0")
	}

	//假关机开机
	if _, ok := result["fake_online"]; ok {
		info := map[string]interface{}{
			"comm_time": time.Now(),
			"state":     "3",
		}
		storage.SetRunInfo(message.Header.Imei, info)
		if err != nil {
			log.Warnf("%v update state failed %v", session.ID(), err)
		}
		log.Warnf("%v update fake online", session.ID())
		storage.Rdb.HSet(context.Background(), fmt.Sprintf("imei_%v", session.ID()), "fake_online", "1")
	}

	err = storage.UpdateCmdResponse(session.ID(), timeid, entity.Content)
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handle0116(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0116)
	session.UserData["short_record"] = ShortRecord{
		Imei:      session.ID(),
		Writer:    common.NewWriter(),
		StartTime: time.Now(),
		Schedule:  0.0,
	}
	log.Infof("%v handle 0116 %v", session.ID(), entity)
	result, err := storage.GetCmdLog(session.ID(), 10, session.Protocol)
	if err != nil {
		return
	}
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}

	var content string
	if entity.RecordStatus == 0 {
		content = "设置成功"
	} else {
		content = "设置失败"
	}
	err = storage.UpdateCmdResponse(session.ID(), timeid, content)
	if err != nil {
		log.Infof("err %v", err)
	}

	info := map[string]interface{}{
		"record_state": entity.RecordStatus,
	}
	storage.SetRunInfo(message.Header.Imei, info)
	storage.SetRecordSchedule(session.ID(), 0)
}

func handle0117(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0117)

	/*
	   这里Encode的时候，会把packet的内容写走
	   log.Infof("body is %v", entity)
	   d, _ := message.Encode()
	   log.Infof("0117 raw msg %x", common.GetHex(d))
	   bdata, _ := entity.Encode()
	   log.Infof("0117 raw body %x", common.GetHex(bdata))
	*/

	shortRecord, ok := session.UserData["short_record"].(ShortRecord)
	if ok {
		shortRecord.Schedule = float32(entity.PkgNo) * 100.0 / float32(entity.PkgSize)
		storage.SetRecordSchedule(session.ID(), shortRecord.Schedule)
		buffer, _ := ioutil.ReadAll(entity.Packet)
		shortRecord.Writer.Write(buffer)
	}

	if entity.PkgNo == entity.PkgSize {
		err := storage.CheckAsValue(shortRecord.Imei, "2")
		if err != nil { //没有增值服务
			log.Warnf("%v no record time left", shortRecord.Imei)
			session.ReplyShortRecord(entity.PkgNo)
			return
		}

		fileSize := len(shortRecord.Writer.Bytes())
		duration := calDuration(fileSize)
		fileName := fmt.Sprintf("%v_%v.amr", shortRecord.Imei, shortRecord.StartTime.Unix())
		reader := bytes.NewReader(shortRecord.Writer.Bytes())
		err = storage.UploadFile("record", fileName, reader, int64(fileSize))
		if err != nil {
			log.Warnf("upload record failed %v", err)
		}

		record := storage.Record{
			Imei:     shortRecord.Imei,
			Time:     shortRecord.StartTime.Unix(),
			Duration: duration,
			Filename: fileName,
			Status:   0,
		}
		err = storage.InsertRecord(record)
		if err != nil {
			log.Warnf("insert record failed %v", err)
		}

		storage.UseAsValue(shortRecord.Imei, "2", duration)
		shortRecord.Writer.Reset()
	}

	data, _ := message.Encode()
	log.Infof("0117 raw msg %x", common.GetHex(data))
	session.ReplyShortRecord(entity.PkgNo)
}

func handle0109(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0109)
	log.Infof("handle 0109 %v", entity)
	session.ReplyTime()
}

func handle0003(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0003)
	log.Infof("handle 0109 %v", entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0105(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0105)
	session.Reply8125()
	log.Infof("%v handle 0105 %v", message.Header.Imei, entity)
}

func handle0108(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0108)
	log.Infof("%v handle 0108 %v", message.Header.Imei, entity)
	session.Reply8108()
}

func handle0210(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0210)
	log.Infof("%v handle 0210 %v", message.Header.Imei, entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0115(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0115)
	log.Infof("%v, handle 0115 %v", message.Header.Imei, entity)
	delete(session.UserData, "short_record")
	storage.DelRunInfoFields(session.ID(), []string{"record_state"})
	//session.Reply8115(entity.SessionId)
}

func handle0120(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0120)
	log.Infof("%v handle 0120 %v", message.Header.Imei, entity)
	session.UserData["vor_record"] = &VorRecord{
		Imei:        session.ID(),
		Writer:      common.NewWriter(),
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		FirstPacket: true,
		PkgCnt:      0,
	}

	info := map[string]interface{}{
		"record_state": 3,
	}
	storage.SetRunInfo(message.Header.Imei, info)

	storage.SetVorSwitch(session.ID(), 1)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0118(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0118)

	currBeginTime := entity.Time

	vorRecord, ok := session.UserData["vor_record"].(*VorRecord)
	if !ok {
		session.UserData["vor_record"] = &VorRecord{
			Imei:        session.ID(),
			Writer:      common.NewWriter(),
			StartTime:   time.Now(),
			EndTime:     time.Now(),
			FirstPacket: true,
			PkgCnt:      0,
		}
		vorRecord, _ = session.UserData["vor_record"].(*VorRecord)
	}

	if vorRecord.FirstPacket {
		vorRecord.StartTime = entity.Time
		vorRecord.EndTime = vorRecord.StartTime.Add(time.Second * 10)
		vorRecord.FirstPacket = false
	}
	if entity.PkgNo == 1 { //组包
		if currBeginTime.Sub(vorRecord.EndTime).Seconds() < 3 &&
			currBeginTime.Sub(vorRecord.StartTime).Seconds() < 59 &&
			vorRecord.PkgCnt < 47 {
			vorRecord.EndTime = currBeginTime.Add(time.Second * 10)
			buffer, err := ioutil.ReadAll(entity.Packet)
			log.Infof("%v rayjay buffer len %v, err is %v", message.Header.Imei, len(buffer), err)
			vorRecord.Writer.Write(buffer)
			vorRecord.PkgCnt += 1
		} else { //上报当前缓存录音
			info := map[string]interface{}{
				"comm_time": time.Now(),
				"state":     "3",
			}
			storage.SetRunInfo(message.Header.Imei, info)

			err := storage.CheckAsValue(vorRecord.Imei, "2")
			if err != nil { //没有增值服务
				log.Warnf("%v no record time left", vorRecord.Imei)
				//重新初始化并写入第一个包
				vorRecord.PkgCnt = 0
				vorRecord.FirstPacket = true
				buffer, _ := ioutil.ReadAll(entity.Packet)
				vorRecord.Writer.Write(buffer)
				vorRecord.PkgCnt += 1
				session.ReplyVorRecord(entity)
				return
			}

			fileName := fmt.Sprintf("%v_%v.amr", vorRecord.Imei, vorRecord.StartTime.Unix())
			fileSize := len(vorRecord.Writer.Bytes())
			duration := calDuration(fileSize)
			reader := bytes.NewReader(vorRecord.Writer.Bytes())
			err = storage.UploadFile("record", fileName, reader, int64(fileSize))
			if err != nil {
				log.Warnf("upload record failed %v", err)
			}
			record := storage.Record{
				Imei:     vorRecord.Imei,
				Time:     vorRecord.StartTime.Unix(),
				Duration: duration,
				Filename: fileName,
				Status:   0,
			}
			err = storage.InsertRecord(record)
			if err != nil {
				log.Warnf("insert record failed %v", err)
			}
			vorRecord.Writer.Reset()

			storage.UseAsValue(vorRecord.Imei, "2", duration)

			//重新初始化并写入第一个包
			vorRecord.PkgCnt = 0
			vorRecord.FirstPacket = true
			buffer, _ := ioutil.ReadAll(entity.Packet)
			vorRecord.Writer.Write(buffer)
			vorRecord.PkgCnt += 1
		}
	} else {
		buffer, _ := ioutil.ReadAll(entity.Packet)
		vorRecord.Writer.Write(buffer)
		vorRecord.PkgCnt += 1
	}
	log.Infof("%v vor %v", message.Header.Imei, *vorRecord)

	session.ReplyVorRecord(entity)
}

func handle0119(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0119)
	log.Infof("%v handle 0119 %v", message.Header.Imei, entity)
	storage.DelRunInfoFields(session.ID(), []string{"record_state"})
	storage.SetVorSwitch(session.ID(), 0)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0001(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0001)
	log.Infof("%v handle 0001 %v", session.ID(), entity)

	if session.Protocol == 5 && entity.ReplyMsgID == uint16(jt808.MsgT808_0x8300) {
		return
	}

	if session.Protocol == 7 && entity.ReplyMsgID == uint16(jt808.MsgT808_0x8300) {
		return
	}

	result, err := storage.GetCmdLog(session.ID(), entity.ReplyMsgSerialNo, session.Protocol)
	if err != nil {
		return
	}
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}

	//同步定位模式
	if mode, ok := result["mode"]; ok {
		err = storage.UpdateMode(session.ID(), mode)
		if err != nil {
			log.Warnf("%v update mode failed %v", session.ID(), err)
		}
	}

	//同步震动报警阈值
	if value, ok := result["shake_value"]; ok {
		shakeValue, _ := strconv.Atoi(value)
		err = storage.UpdateShakeValue(session.ID(), shakeValue)
		if err != nil {
			log.Warnf("%v update mode failed %v", session.ID(), err)
		}
	}

	var content string
	if entity.Result == 0 {
		content = "设置成功"
	} else {
		content = "设置失败"
	}
	err = storage.UpdateCmdResponse(session.ID(), timeid, content)
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handle0107(session *jt808_base.Session, message *jt808.Message) {
	//除了在这儿之外，其他都通过1107处理
	/*
		if session.Protocol != 7 {
			session.Reply(message, jt808.T808_0x8100_ResultSuccess)
			return
		}
	*/

	entity := message.Body.(*jt808.T808_0x0107)
	log.Infof("%v handle 0107 %v", session.ID(), entity)

	iccid := entity.Iccid

	fTrans := func(s *string) {
		runes := []rune(*s) // 将字符串转换为rune切片，处理Unicode字符
		for i := 0; i < len(runes); i++ {
			switch runes[i] {
			case 0x3A:
				runes[i] = 0x41
			case 0x3B:
				runes[i] = 0x42
			case 0x3C:
				runes[i] = 0x43
			case 0x3D:
				runes[i] = 0x44
			case 0x3E:
				runes[i] = 0x45
			case 0x3F:
				runes[i] = 0x46
			}

			if runes[i] > 0x60 && runes[i] < 0x7B {
				runes[i] = runes[i] ^ 32
			}
		}
		*s = string(runes) // 将rune切片转换回字符串
	}
	fTrans(&iccid)

	deviceInfo, err := storage.GetDevice(session.ID())
	if err != nil {
		session.Reply(message, jt808.T808_0x8100_ResultSuccess)
		return
	}
	if v, ok := deviceInfo["iccid"]; ok {
		if v != iccid {
			log.Infof("%v update iccid from %v to %v", session.ID(), iccid, entity.Iccid)
			err = storage.UpdateIccid(session.ID(), entity.Iccid)
			if err != nil {
				log.Warnf("%v update iccid failed %v", session.ID(), err)
			}
		}
	}

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0112(session *jt808_base.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0112)
	log.Infof("%v handle 0112 %v", session.ID(), entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1006(session *jt808_base.Session, message *jt808.Message) {
	//entity := message.Body.(*jt808.T808_0x1006)
	//log.Infof("%v handle 1006 %v", session.ID(), entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}
