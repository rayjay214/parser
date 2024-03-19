package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rayjay214/parser/common"
	"github.com/rayjay214/parser/jt808"
	"github.com/rayjay214/parser/jt808/extra"
	"io/ioutil"
	"net/http"
	"strings"
)

type Bts struct {
	Lac    uint16
	Cellid uint32
	Rssi   uint8
}

type LbsInfo struct {
	Mcc     uint16
	Mnc     uint8
	BtsList []Bts
}

type Mac struct {
	MacAddr string
	Rssi    uint8
}

type WifiInfo struct {
	MacList []Mac
}

type LbsResp struct {
	Errorcode int     `json:"errorcode"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	LocType   int
}

func getLbsLocation(entity *jt808.T808_0x0200, lbsResp *LbsResp) error {
	url := "http://114.215.191.234/locapi"

	var lbsInfo LbsInfo
	var wifiInfo WifiInfo

	fParseE1 := func(lbsContent []byte, info *LbsInfo) error {
		reader := common.NewReader(lbsContent)

		var err error
		info.Mcc, err = reader.ReadUint16()
		if err != nil {
			return err
		}
		Mnc, err := reader.ReadUint16()
		if err != nil {
			return err
		}
		info.Mnc = uint8(Mnc)

		btsNum := (len(lbsContent) - 4) / 8
		for i := 0; i < btsNum; i++ {
			var bts Bts
			reader.ReadByte() //useless byte
			bts.Lac, err = reader.ReadUint16()
			if err != nil {
				return err
			}
			bts.Cellid, err = reader.ReadUint32()
			if err != nil {
				return err
			}
			rssi, err := reader.ReadByte()
			if err != nil {
				return err
			}
			bts.Rssi = rssi*2 - 113
			info.BtsList = append(info.BtsList, bts)
		}
		return nil
	}

	fParseEB := func(lbsContent []byte, info *LbsInfo) error {
		reader := common.NewReader(lbsContent)

		var err error
		info.Mcc, err = reader.ReadUint16()
		if err != nil {
			return err
		}
		info.Mnc, err = reader.ReadByte()
		if err != nil {
			return err
		}

		btsNum := (len(lbsContent) - 3) / 5
		for i := 0; i < btsNum; i++ {
			var bts Bts
			bts.Lac, err = reader.ReadUint16()
			if err != nil {
				return err
			}
			cellid, err := reader.ReadUint16()
			if err != nil {
				return err
			}
			bts.Cellid = uint32(cellid)
			rssi, err := reader.ReadByte()
			if err != nil {
				return err
			}
			bts.Rssi = rssi*2 - 113
			info.BtsList = append(info.BtsList, bts)
		}
		return nil
	}

	fParseEE := func(lbsContent []byte, info *LbsInfo) error {
		reader := common.NewReader(lbsContent)

		var err error
		info.Mcc, err = reader.ReadUint16()
		if err != nil {
			return err
		}
		info.Mnc, err = reader.ReadByte()
		if err != nil {
			return err
		}

		var bts Bts

		bts.Lac, err = reader.ReadUint16()
		if err != nil {
			return err
		}
		bts.Cellid, err = reader.ReadUint32()
		if err != nil {
			return err
		}
		rssi, err := reader.ReadByte()
		if err != nil {
			return err
		}
		bts.Rssi = rssi*2 - 113
		info.BtsList = append(info.BtsList, bts)

		return nil
	}

	fParseEC := func(wifiContent []byte, info *WifiInfo) error {
		reader := common.NewReader(wifiContent)

		wifiNum := len(wifiContent) / 7

		var err error
		for i := 0; i < wifiNum; i++ {
			var mac Mac
			var byteMac []byte
			byteMac, err = reader.Read(6)
			if err != nil {
				return err
			}
			strMac := hex.EncodeToString(byteMac)
			var parts []string
			for i := 0; i < len(strMac); i += 2 {
				end := i + 2
				if end > len(strMac) {
					end = len(strMac)
				}
				parts = append(parts, strMac[i:end])
			}
			mac.MacAddr = strings.Join(parts, ":")

			mac.Rssi, err = reader.ReadByte()
			if err != nil {
				return err
			}
			info.MacList = append(info.MacList, mac)
		}
		return nil
	}

	for _, ext := range entity.Extras {
		switch ext.ID() {
		case extra.Extra_0xeb{}.ID():
			lbsContent := ext.(*extra.Extra_0xeb).Data()
			fParseEB(lbsContent, &lbsInfo)
		case extra.Extra_0xe1{}.ID():
			lbsContent := ext.(*extra.Extra_0xe1).Data()
			fParseE1(lbsContent, &lbsInfo)
		case extra.Extra_0xee{}.ID():
			lbsContent := ext.(*extra.Extra_0xee).Data()
			fParseEE(lbsContent, &lbsInfo)
		case extra.Extra_0xec{}.ID():
			wifiContent := ext.(*extra.Extra_0xec).Data()
			fParseEC(wifiContent, &wifiInfo)
		}
	}

	if len(wifiInfo.MacList) == 0 && len(lbsInfo.BtsList) == 0 {
		return errors.New("no valid lbs or wifi info")
	}

	if len(lbsInfo.BtsList) > 0 {
		lbsResp.LocType = 1
	}

	if len(wifiInfo.MacList) > 0 {
		lbsResp.LocType = 2
	}

	body := make(map[string]interface{})

	if len(lbsInfo.BtsList) > 0 {
		body["accesstype"] = 0
		body["bts"] = fmt.Sprintf("%d,%d,%d,%d,%d", lbsInfo.Mcc, lbsInfo.Mnc, lbsInfo.BtsList[0].Lac, lbsInfo.BtsList[0].Cellid, lbsInfo.BtsList[0].Rssi)
		var btsList []string
		for _, bts := range lbsInfo.BtsList {
			strBts := fmt.Sprintf("%d,%d,%d,%d,%d", lbsInfo.Mcc, lbsInfo.Mnc, bts.Lac, bts.Cellid, bts.Rssi)
			btsList = append(btsList, strBts)
		}
		body["nearbts"] = strings.Join(btsList, "|")
	}

	if len(wifiInfo.MacList) > 0 {
		body["accesstype"] = 1
		var macList []string
		for _, mac := range wifiInfo.MacList {
			strMac := fmt.Sprintf("%s,%d", mac.MacAddr, mac.Rssi)
			macList = append(macList, strMac)
		}
		body["macs"] = strings.Join(macList, "|")
	}
	byteData, _ := json.Marshal(body)
	//log.Infof("post data is %v", string(byteData))
	reader := bytes.NewReader(byteData)

	request, err := http.NewRequest("POST", url, reader)
	defer request.Body.Close()
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//log.Infof("resp is %v", string(respBytes))

	err = json.Unmarshal(respBytes, lbsResp)
	if err != nil {
		return err
	}

	return nil

}

func calDuration(fileSize int) int {
	quotient := fileSize / 702
	remainder := fileSize % 702

	// 如果余数大于等于除数的一半，向上取整
	if remainder >= 702/2 {
		quotient++
	}
	return quotient
}
