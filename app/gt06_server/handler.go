package main

import (
	"fmt"
	"github.com/rayjay214/parser/protocol/gt06"
	"github.com/rayjay214/parser/server_base/gt06_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func handle01(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x01)
	fmt.Printf("%v:handle 01 %v, %v\n", session.ID(), message, entity)

	set, err := storage.SetStartTime(session.ID())
	if err == nil && set {
		storage.UpdateStartTime(session.ID())
	}

	info := map[string]interface{}{
		"state":     "4",
		"comm_time": time.Now(),
	}
	_ = storage.SetRunInfo(session.ID(), info)

	session.CommonReply(entity.Proto)
}

func handle12(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x12)
	fmt.Printf("%v:handle 12 %v, %v\n", session.ID(), message, entity)
}

func handleA1(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0xa1)

	var lbsResp LbsResp
	var lbsInfo LbsInfo
	lbsInfo.Mcc = entity.Mcc
	lbsInfo.Mnc = entity.Mnc
	for _, item := range entity.LbsInfoList {
		var bts Bts
		if item.Lac != 0 {
			bts.Cellid = uint32(item.CellId)
			bts.Lac = item.Lac
			bts.Rssi = item.Rssi
			lbsInfo.BtsList = append(lbsInfo.BtsList, bts)
		}
	}
	getLbsLocation(lbsInfo, &lbsResp)

	date := entity.Time.Format("20060102")
	iDate, _ := strconv.Atoi(date)

	loc := storage.Location{
		Imei:      session.ID(),
		Date:      iDate,
		Time:      entity.Time.Unix(),
		Direction: 0,
		Lat:       int64(lbsResp.Lat * 1000000),
		Lng:       int64(lbsResp.Lng * 1000000),
		Speed:     0,
		Type:      7,
		Wgs:       "",
	}
	err := storage.InsertLocation(loc)
	if err != nil {
		log.Warnf("insert location err %v", err)
	}

	info := map[string]interface{}{
		"comm_time": time.Now(),
	}
	info["lat"] = lbsResp.Lat
	info["lng"] = lbsResp.Lng
	info["loc_type"] = 7
	info["loc_time"] = entity.Time

	storage.SetRunInfo(session.ID(), info)

	fmt.Printf("handle a1 %v, %v\n", message, entity)
}

func handle94(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x94)
	fmt.Printf("%v:handle 94 %v, %v\n", session.ID(), message, entity)
	if entity.SubProto == 0x0a {
		var hexStrings []string
		for _, b := range entity.Content[16:26] {
			hexStrings = append(hexStrings, fmt.Sprintf("%02X", b))
		}
		iccid := strings.Join(hexStrings, "")
		fmt.Println(iccid)
		deviceInfo, err := storage.GetDevice(session.ID())
		if err != nil {
			return
		}
		if lastIccid, ok := deviceInfo["iccid"]; ok {
			if lastIccid != iccid {
				log.Infof("%v update iccid from %v to %v", session.ID(), lastIccid, iccid)
				err = storage.UpdateIccid(session.ID(), iccid)
				if err != nil {
					log.Warnf("%v update iccid failed %v", session.ID(), err)
				}
			}
		}
	}
}

func handle20(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x20)

	var lbsResp LbsResp
	var lbsInfo LbsInfo
	lbsInfo.Mcc = entity.Mcc
	lbsInfo.Mnc = uint16(entity.Mnc)
	for _, item := range entity.LbsInfoList {
		var bts Bts
		if item.Lac != 0 {
			bts.Cellid = uint32(item.CellId)
			bts.Lac = uint32(item.Lac)
			bts.Rssi = item.Rssi
			lbsInfo.BtsList = append(lbsInfo.BtsList, bts)
		}
	}

	//忽略wifi

	getLbsLocation(lbsInfo, &lbsResp)

	//07设备定位时间无用，采用服务器的
	if session.Protocol == 3 {
		entity.Time = time.Now()
	}

	date := entity.Time.Format("20060102")
	iDate, _ := strconv.Atoi(date)

	loc := storage.Location{
		Imei:      session.ID(),
		Date:      iDate,
		Time:      entity.Time.Unix(),
		Direction: 0,
		Lat:       int64(lbsResp.Lat * 1000000),
		Lng:       int64(lbsResp.Lng * 1000000),
		Speed:     0,
		Type:      7,
		Wgs:       "",
	}

	handleLocation(session.ID(), loc, lbsResp, entity.Time)
}

func handle13(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x13)
	fmt.Printf("%v:handle 13 %v, %v\n", session.ID(), message, entity)

	info := map[string]interface{}{
		"state":     "4",
		"comm_time": time.Now(),
		"power":     entity.Voltage,
	}
	_ = storage.SetRunInfo(session.ID(), info)

	session.CommonReply(entity.Proto)
}

func handle16(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x16)

	var lbsResp LbsResp
	var lbsInfo LbsInfo
	lbsInfo.Mcc = entity.Mcc
	lbsInfo.Mnc = uint16(entity.Mnc)

	var bts Bts
	bts.Cellid = uint32(entity.CellId)
	bts.Lac = uint32(entity.Lac)
	bts.Rssi = 0
	lbsInfo.BtsList = append(lbsInfo.BtsList, bts)

	getLbsLocation(lbsInfo, &lbsResp)

	//07设备定位时间无用，采用服务器的
	if session.Protocol == 3 {
		entity.Time = time.Now()
	}

	date := entity.Time.Format("20060102")
	iDate, _ := strconv.Atoi(date)

	loc := storage.Location{
		Imei:      session.ID(),
		Date:      iDate,
		Time:      entity.Time.Unix(),
		Direction: 0,
		Lat:       int64(lbsResp.Lat * 1000000),
		Lng:       int64(lbsResp.Lng * 1000000),
		Speed:     0,
		Type:      7,
		Wgs:       "",
	}

	handleLocation(session.ID(), loc, lbsResp, entity.Time)

	info := map[string]interface{}{
		"power": entity.Voltage,
	}
	_ = storage.SetRunInfo(session.ID(), info)

	session.CommonReply(entity.Proto)
}

func handleLocation(imei uint64, loc storage.Location, lbsResp LbsResp, locTime time.Time) {
	if loc.Lat == 0 || loc.Lng == 0 {
		return
	}

	err := storage.InsertLocation(loc)
	if err != nil {
		log.Warnf("insert location err %v", err)
	}

	info := map[string]interface{}{
		"comm_time": time.Now(),
		"state":     "4",
	}
	info["lat"] = lbsResp.Lat
	info["lng"] = lbsResp.Lng
	info["loc_type"] = 7
	info["loc_time"] = locTime

	storage.SetRunInfo(imei, info)
}

func handle15(session *gt06_base.Session, message *gt06.Message) {
	entity := message.Body.(*gt06.Kks_0x15)
	fmt.Printf("%v:handle 15 %v, %v\n", session.ID(), message, entity)

	result, err := storage.GetCmdLog(session.ID(), uint16(entity.SysFlag))
	if err != nil {
		return
	}
	var timeid uint64
	if v, ok := result["timeid"]; ok {
		timeid, _ = strconv.ParseUint(v, 10, 64)
	}

	err = storage.UpdateCmdResponse(session.ID(), timeid, entity.Content)
}
