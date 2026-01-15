package main

import (
	"context"
	"fmt"
	"github.com/rayjay214/parser/protocol/hl3g"
	"github.com/rayjay214/parser/server_base/hl3g_base"
	"github.com/rayjay214/parser/storage"
	"github.com/rayjay214/parser/util"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func handleLK2(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_LK2)
	log.Infof("%v:handle lk2 %v, %v", session.ID(), message, entity)

	fakeOnline, _ := storage.Rdb.HGet(context.Background(), fmt.Sprintf("imei_%v", session.ID()), "fake_online_state").Result()
	if fakeOnline == "0" {
		session.CommonReply(message.Header.Imei, message.Header.Proto)
		return
	}

	set, err := storage.SetStartTime(session.ID())
	if err == nil && set {
		storage.UpdateStartTime(session.ID())
	}

	info := map[string]interface{}{
		"state":     "3",
		"comm_time": time.Now(),
	}
	_ = storage.SetRunInfo(session.ID(), info)

	session.CommonReply(message.Header.Imei, message.Header.Proto)
}

func handleCCID(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_CCID)
	log.Infof("%v:handle ccid %v, %v", session.ID(), message, entity)

	log.Infof("%v update iccid to %v", session.ID(), entity.Iccid)
	err := storage.UpdateIccid(session.ID(), entity.Iccid)
	if err != nil {
		log.Warnf("%v update iccid failed %v", session.ID(), err)
	}

	session.CommonReply(message.Header.Imei, message.Header.Proto)
}

func handleGS1(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_GS1)
	log.Infof("%v:handle gs1 %v, %v", session.ID(), message, entity)

	lbsInfo := hl3g.LbsInfo{
		Mcc:    entity.Mcc,
		Mnc:    entity.Mnc,
		Lac:    entity.Lac,
		CellId: entity.CellId,
		Rssi:   entity.Rssi,
	}
	var infoList []hl3g.LbsInfo
	infoList = append(infoList, lbsInfo)

	var resp LbsResp
	getLbsLocation(infoList, &resp, session.ID())

	log.Infof("lbs resp %v", resp)
	nowStr := time.Now().Format("20060102150405")
	latStr := fmt.Sprintf("%v", resp.Lat)
	lngStr := fmt.Sprintf("%v", resp.Lng)

	session.Gs1Reply(message.Header.Imei, "GS", latStr, lngStr, nowStr)
}

func handleUD(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_UD)
	log.Infof("%v:handle ud %v, %v", session.ID(), message, entity)
	handleLocation(session.ID(), &entity.LocInfo, session.Protocol)
}

func handleUD2(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_UD2)
	log.Infof("%v:handle ud2 %v, %v", session.ID(), message, entity)
	handleLocation(session.ID(), &entity.LocInfo, session.Protocol)
}

func handleAL(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_AL)
	log.Infof("%v:handle al %v, %v", session.ID(), message, entity)
	handleLocation(session.ID(), &entity.LocInfo, session.Protocol)

	session.CommonReply(message.Header.Imei, message.Header.Proto)
}

func handleSTU(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_STU)
	log.Infof("%v:handle stu %v, %v", session.ID(), message, entity)
}

func handleRESET(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_STU)
	log.Infof("%v:handle stu %v, %v", session.ID(), message, entity)
}

func handleFACTORY(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_STU)
	log.Infof("%v:handle stu %v, %v", session.ID(), message, entity)
}

func handleVERNO(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_VERNO)
	log.Infof("%v:handle verno %v, %v", session.ID(), message, entity)

	result, err := storage.GetCmdLog(session.ID(), 0, session.Protocol)
	if err != nil {
		return
	}
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}
	err = storage.UpdateCmdResponse(session.ID(), timeid, entity.Version)
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handleTC(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_TC)
	log.Infof("%v:handle tc %v, %v", session.ID(), message, entity)

	result, err := storage.GetCmdLog(session.ID(), 0, session.Protocol)
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

func handleMODEWORK(session *hl3g_base.Session, message *hl3g.Message) {
	log.Infof("%v:handle modework %v", session.ID(), message)
	result, err := storage.GetCmdLog(session.ID(), 0, session.Protocol)
	if err != nil {
		return
	}
	log.Infof("%v:result %v", session.ID(), result)
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

	err = storage.UpdateCmdResponse(session.ID(), timeid, "OK")
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handleUPLOAD(session *hl3g_base.Session, message *hl3g.Message) {
	log.Infof("%v:handle upload %v", session.ID(), message)
	result, err := storage.GetCmdLog(session.ID(), 0, session.Protocol)
	if err != nil {
		return
	}
	log.Infof("%v:result %v", session.ID(), result)
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

	err = storage.UpdateCmdResponse(session.ID(), timeid, "OK")
	if err != nil {
		log.Infof("err %v", err)
	}
}

func handleCLOSEMODE(session *hl3g_base.Session, message *hl3g.Message) {
	log.Infof("%v:handle closemode %v", session.ID(), message)
	result, err := storage.GetCmdLog(session.ID(), 0, session.Protocol)
	if err != nil {
		return
	}
	log.Infof("%v:result %v", session.ID(), result)
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}

	storage.SetFakeOnlineState(result, session.ID())

	err = storage.UpdateCmdResponse(session.ID(), timeid, "OK")
	if err != nil {
		log.Infof("err %v", err)
	}
}

type LocSource int

const (
	LocGPS  LocSource = 0
	LocLBS  LocSource = 1
	LocWiFi LocSource = 2
)

func calcLocType(source LocSource, speed int) int {
	base := int(source)
	if speed > 0 {
		return base // 0 / 1 / 2
	}
	return base + 3 // 3 / 4 / 5
}

func handleLocation(imei uint64, info *hl3g.LocationInfo, protocol int) {
	log.Infof("%v:location info %v", imei, info)

	now := time.Now()
	zone, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("020106150405", info.Date+info.Time, zone)
	if (now.Unix() - t.Unix()) > 7*3600 {
		t = t.Add(time.Hour * 8)
	}

	/*
		zone, _ := time.LoadLocation("Asia/Shanghai")
		tUtc, _ := time.Parse("020106150405", info.Date+info.Time)
		var tShanghai time.Time
		if info.Located == "A" {
			tShanghai = tUtc.In(zone)
		} else {
			tShanghai, _ = time.ParseInLocation("020106150405", info.Date+info.Time, zone)
		}

	*/

	iSpeed, _ := strconv.Atoi(info.Speed)
	iDirection, _ := strconv.Atoi(info.Direction)

	date := t.Format("20060102")
	iDate, _ := strconv.Atoi(date)

	runinfo := map[string]interface{}{
		"comm_time": time.Now(),
		"power":     info.Power,
		"signal":    info.Gsm,
		"satellite": info.Power,
		"loc_time":  t,
	}

	log.Infof("loc_time is %v", t)

	if iSpeed > 0 {
		runinfo["state"] = 2
	} else {
		runinfo["state"] = 3
	}

	loc := storage.Location{
		Imei:      imei,
		Date:      iDate,
		Time:      t.Unix(),
		Direction: uint16(iDirection),
		Speed:     uint16(iSpeed),
		Wgs:       "",
	}

	var (
		source LocSource
		fLat   float64
		fLng   float64
	)

	if info.Located == "A" {
		source = LocGPS
		fLat, _ = strconv.ParseFloat(info.Lat, 64)
		fLng, _ = strconv.ParseFloat(info.Lng, 64)

	} else {
		var resp LbsResp

		if len(info.Wifi) > 0 {
			getWifiLocation(info.Wifi, &resp, imei)
			source = LocWiFi

			if resp.Lat == 0 && resp.Lng == 0 && len(info.Lbs) > 0 {
				getLbsLocation(info.Lbs, &resp, imei)
				source = LocLBS
			}
		} else {
			getLbsLocation(info.Lbs, &resp, imei)
			source = LocLBS
		}

		fLat = resp.Lat
		fLng = resp.Lng
	}

	loc.Lat = int64(fLat * 1e6)
	loc.Lng = int64(fLng * 1e6)
	loc.Type = calcLocType(source, iSpeed)

	runinfo["lat"] = fLat
	runinfo["lng"] = fLng
	runinfo["loc_type"] = loc.Type

	if loc.Lat == 0 || loc.Lng == 0 {
		return
	}

	if err := storage.InsertLocation(loc); err != nil {
		log.Warnf("insert location err %v", err)
	}

	storage.SetRunInfo(imei, runinfo)

	util.CheckFence(&loc)

	checkAlarm(info, &loc)
}

func checkAlarm(info *hl3g.LocationInfo, loc *storage.Location) {
	alarmValue, _ := strconv.ParseUint(info.Status, 16, 32)

	alarm := storage.Alarm{
		Imei:  loc.Imei,
		Time:  time.Unix(loc.Time, 0),
		Lat:   loc.Lat,
		Lng:   loc.Lng,
		Speed: loc.Speed,
	}

	const VibrateBit = 4
	const SOSBit = 16
	const LowpowerBit = 17

	if (alarmValue>>VibrateBit)&1 == 1 {
		alarm.Type = "1"
		storage.InsertAlarm(alarm)
	}
	if (alarmValue>>SOSBit)&1 == 1 {
		alarm.Type = "4"
		storage.InsertAlarm(alarm)
	}
	if (alarmValue>>LowpowerBit)&1 == 1 {
		alarm.Type = "3"
		storage.InsertAlarm(alarm)
	}
}
