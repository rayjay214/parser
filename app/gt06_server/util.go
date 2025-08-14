package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Bts struct {
	Lac    uint32
	Cellid uint32
	Rssi   uint8
}

type LbsInfo struct {
	Mcc     uint16
	Mnc     uint16
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

func getLbsLocation(lbsInfo LbsInfo, lbsResp *LbsResp, imei uint64) error {
	//url := "http://114.215.191.234/locapi"
	url := "http://121.196.220.14/locapi"

	lbsResp.LocType = 1

	body := make(map[string]interface{})
	body["accesstype"] = 0
	body["bts"] = fmt.Sprintf("%d,%d,%d,%d,%d", lbsInfo.Mcc, lbsInfo.Mnc, lbsInfo.BtsList[0].Lac, lbsInfo.BtsList[0].Cellid, lbsInfo.BtsList[0].Rssi)
	var btsList []string
	for _, bts := range lbsInfo.BtsList {
		strBts := fmt.Sprintf("%d,%d,%d,%d,%d", lbsInfo.Mcc, lbsInfo.Mnc, bts.Lac, bts.Cellid, bts.Rssi)
		btsList = append(btsList, strBts)
	}
	body["nearbts"] = strings.Join(btsList, "|")
	byteData, _ := json.Marshal(body)
	fmt.Printf("%v post data is %v\n", imei, string(byteData))
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
	fmt.Printf("%v resp is %v\n", imei, string(respBytes))

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
