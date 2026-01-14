package main

import (
	"fmt"
	"github.com/rayjay214/parser/protocol/hl3g"
	"github.com/rayjay214/parser/server_base/hl3g_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func handleLK2(session *hl3g_base.Session, message *hl3g.Message) {
	entity := message.Body.(*hl3g.HL3G_LK2)
	log.Infof("%v:handle lk2 %v, %v", session.ID(), message, entity)

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

	zone, _ := time.LoadLocation("Asia/Shanghai")
	tUtc, _ := time.Parse("020106150405", info.Date+info.Time)
	tShanghai := tUtc.In(zone)

	iSpeed, _ := strconv.Atoi(info.Speed)
	iDirection, _ := strconv.Atoi(info.Direction)

	date := tShanghai.Format("20060102")
	iDate, _ := strconv.Atoi(date)

	runinfo := map[string]interface{}{
		"comm_time": time.Now(),
		"power":     info.Power,
		"signal":    info.Gsm,
		"satellite": info.Power,
		"loc_time":  tShanghai,
	}

	log.Infof("loc_time is %v", tShanghai)

	if iSpeed > 0 {
		runinfo["state"] = 2
	} else {
		runinfo["state"] = 3
	}

	loc := storage.Location{
		Imei:      imei,
		Date:      iDate,
		Time:      tShanghai.Unix(),
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
}
