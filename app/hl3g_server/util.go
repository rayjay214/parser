package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/hl3g"
	"io/ioutil"
	"net/http"
	"strings"
)

type Bts struct {
	Lac    string
	Cellid string
	Rssi   string
}

type LbsInfo struct {
	Mcc     string
	Mnc     string
	BtsList []Bts
}

type Mac struct {
	MacAddr string
	Rssi    string
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

func getLbsLocation(lbsInfo []hl3g.LbsInfo, lbsResp *LbsResp, imei uint64) error {
	url := "http://121.196.220.14/locapi"

	lbsResp.LocType = 1

	body := make(map[string]interface{})
	body["accesstype"] = 0
	body["bts"] = fmt.Sprintf("%s,%s,%s,%s,%s", lbsInfo[0].Mcc, lbsInfo[0].Mnc, lbsInfo[0].Lac, lbsInfo[0].CellId, lbsInfo[0].Rssi)
	var btsList []string
	for _, bts := range lbsInfo {
		strBts := fmt.Sprintf("%s,%s,%s,%s,%s", bts.Mcc, bts.Mnc, bts.Lac, bts.CellId, bts.Rssi)
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

func getWifiLocation(wifiInfo []hl3g.WifiInfo, lbsResp *LbsResp, imei uint64) error {
	url := "http://121.196.220.14/locapi"

	lbsResp.LocType = 1

	body := make(map[string]interface{})
	body["accesstype"] = 1
	var macList []string
	for _, wifi := range wifiInfo {
		strMacs := fmt.Sprintf("%s,%s", wifi.Mac, wifi.Rssi)
		macList = append(macList, strMacs)
	}
	body["macs"] = strings.Join(macList, "|")
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
