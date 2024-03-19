package main

import (
	"bytes"
	"fmt"
	"github.com/rayjay214/parser/common"
	"github.com/rayjay214/parser/jt808"
	"github.com/rayjay214/parser/jt808/extra"
	"github.com/rayjay214/parser/server"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"time"
)

func handle0100(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0100)
	log.Infof("handle 0100 %v", entity)

	session.ReplyRegister(message)
}

func handle0102(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0102)
	log.Infof("handle 0102 %v", entity)

	set, err := storage.SetStartTime(session.ID())
	if err == nil && set {
		storage.UpdateStartTime(session.ID())
	}

	info := map[string]interface{}{
		"comm_time": time.Now(),
	}
	storage.SetRunInfo(message.Header.Imei, info)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0002(session *server.Session, message *jt808.Message) {
	info := map[string]interface{}{
		"comm_time": time.Now(),
	}
	storage.SetRunInfo(message.Header.Imei, info)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0808(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0808)
	session.Protocol = 2 //2013
	log.Infof("handle 0808 %v", entity)
}

// 处理上报位置
func handle0200(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0200)
	//log.Infof("handle 0200 %v", entity)
	handleLocation(message.Header.Imei, entity, session.Protocol)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0704(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0704)
	log.Infof("handle 0704 %v", entity)

	for _, item := range entity.Items {
		handleLocation(message.Header.Imei, &item, session.Protocol)
	}

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0201(session *server.Session, message *jt808.Message) {
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

	for _, ext := range entity.Extras {
		switch ext.ID() {
		case extra.Extra_0x01{}.ID():
			info["distance"] = ext.(*extra.Extra_0x01).Value()
		case extra.Extra_0x04{}.ID():
			v := ext.(*extra.Extra_0x04).Value().(extra.Extra_0xe4_Value)
			info["power"] = v.Power
			info["acc_power"] = v.Status
		case extra.Extra_0xe4{}.ID():
			v := ext.(*extra.Extra_0xe4).Value().(extra.Extra_0xe4_Value)
			info["power"] = v.Power
			info["acc_power"] = v.Status
		case extra.Extra_0x05{}.ID():
			if protocol == 1 {
				if ext.(*extra.Extra_0x05).Value() == 1 {
					info["state"] = "2"
					locTypeBase = 3
				} else {
					info["state"] = "3"
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
						info["state"] = "2"
						locTypeBase = 3
					} else {
						info["state"] = "3"
					}
				} else {
					if entity.Status.GetAccState() {
						info["state"] = "3"
					} else {
						info["state"] = "2"
						locTypeBase = 3
					}
				}
			}
		case extra.Extra_0xf0{}.ID():
			log.Infof("voltage is %v", ext.(*extra.Extra_0xf0).Value())
			info["voltage"] = fmt.Sprintf("%.2fV", float32(ext.(*extra.Extra_0xf0).Value().(uint16)/100))
		case extra.Extra_0xf5{}.ID():
			info["gprs_type"] = ext.(*extra.Extra_0xf5).Value()
		}
	}

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

	} else {
		var lbsResp LbsResp
		err := getLbsLocation(entity, &lbsResp)
		if err != nil {
			return
		}
		if lbsResp.Lng == 0 {
			log.WithFields(log.Fields{
				"imei": imei,
			}).Info("get lbs location failed")
			return
		}

		loc = storage.Location{
			Imei:      imei,
			Date:      iDate,
			Time:      entity.Time.Unix(),
			Direction: entity.Direction,
			Lat:       int64(lbsResp.Lat * 1000000),
			Lng:       int64(lbsResp.Lng * 1000000),
			Speed:     entity.Speed,
			Type:      lbsResp.LocType + locTypeBase,
			Wgs:       "",
		}

		info["lat"] = lbsResp.Lat
		info["lng"] = lbsResp.Lng
		info["loc_type"] = lbsResp.LocType + locTypeBase
		info["loc_time"] = entity.Time

		err = storage.InsertLocation(loc)
		if err != nil {
			log.Warnf("insert location err %v", err)
		}
	}
	storage.SetRunInfo(imei, info)

	checkAlarm(entity, &loc, protocol)
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

func handle1007(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x1007)
	log.Infof("handle 1007 %v", entity)

	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1107(session *server.Session, message *jt808.Message) {
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

func handle1300(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x1300)
	log.Infof("handle 1300 %v", entity)
	result, err := storage.GetCmdLog(session.ID(), entity.AckSeqNo)
	if err != nil {
		return
	}
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}
	err = storage.UpdateCmdResponse(session.ID(), timeid, entity.Content)
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handle0116(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0116)
	session.UserData["short_record"] = ShortRecord{
		Imei:      session.ID(),
		Writer:    common.NewWriter(),
		StartTime: time.Now(),
		Schedule:  0.0,
	}
	log.Infof("%v handle 0116 %v", session.ID(), entity)
	result, err := storage.GetCmdLog(session.ID(), 10)
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

func handle0117(session *server.Session, message *jt808.Message) {
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
		fileSize := len(shortRecord.Writer.Bytes())
		duration := calDuration(fileSize)
		fileName := fmt.Sprintf("%v_%v.amr", shortRecord.Imei, shortRecord.StartTime.Unix())
		reader := bytes.NewReader(shortRecord.Writer.Bytes())
		err := storage.UploadFile("record", fileName, reader, int64(fileSize))
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
		shortRecord.Writer.Reset()
	}

	data, _ := message.Encode()
	log.Infof("0117 raw msg %x", common.GetHex(data))
	session.ReplyShortRecord(entity.PkgNo)
}

func handle0109(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0109)
	log.Infof("handle 0109 %v", entity)
	session.ReplyTime()
}

func handle0003(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0003)
	log.Infof("handle 0109 %v", entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0105(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0105)
	session.Reply8125()
	log.Infof("handle 0105 %v", entity)
}

func handle0108(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0108)
	log.Infof("handle 0108 %v", entity)
	session.Reply8108()
}

func handle0210(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0210)
	log.Infof("handle 0210 %v", entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0115(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0115)
	log.Infof("handle 0115 %v", entity)
	delete(session.UserData, "short_record")
	storage.DelRunInfoFields(session.ID(), []string{"record_state"})
	//session.Reply8115(entity.SessionId)
}

func handle0120(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0120)
	log.Infof("handle 0120 %v", entity)
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

func handle0118(session *server.Session, message *jt808.Message) {
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
			log.Infof("rayjay buffer len %v, err is %v", len(buffer), err)
			vorRecord.Writer.Write(buffer)
			vorRecord.PkgCnt += 1
		} else { //上报当前缓存录音
			fileName := fmt.Sprintf("%v_%v.amr", vorRecord.Imei, vorRecord.StartTime.Unix())
			fileSize := len(vorRecord.Writer.Bytes())
			duration := calDuration(fileSize)
			reader := bytes.NewReader(vorRecord.Writer.Bytes())
			err := storage.UploadFile("record", fileName, reader, int64(fileSize))
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
	log.Infof("vor %v", *vorRecord)

	session.ReplyVorRecord(entity)
}

func handle0119(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0119)
	log.Infof("handle 0119 %v", entity)
	storage.DelRunInfoFields(session.ID(), []string{"record_state"})
	storage.SetVorSwitch(session.ID(), 0)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0001(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0001)
	log.Infof("%v handle 0001 %v", session.ID(), entity)

	result, err := storage.GetCmdLog(session.ID(), entity.ReplyMsgSerialNo)
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

func handle0107(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0107)
	log.Infof("%v handle 0107 %v", session.ID(), entity)

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

func handle0112(session *server.Session, message *jt808.Message) {
	entity := message.Body.(*jt808.T808_0x0112)
	log.Infof("%v handle 0112 %v", session.ID(), entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1006(session *server.Session, message *jt808.Message) {
	//entity := message.Body.(*jt808.T808_0x1006)
	//log.Infof("%v handle 1006 %v", session.ID(), entity)
	session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}
