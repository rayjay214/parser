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
		deviceInfo, err := storage.GetDevice(session.ID())
		if err != nil {
			return
		}
		if lastIccid, ok := deviceInfo["iccid"]; ok {
			if iccid != iccid {
				log.Infof("%v update iccid from %v to %v", session.ID(), lastIccid, iccid)
				err = storage.UpdateIccid(session.ID(), iccid)
				if err != nil {
					log.Warnf("%v update iccid failed %v", session.ID(), err)
				}
			}
		}
	}
}
