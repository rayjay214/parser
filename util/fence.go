package util

import (
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"github.com/pkg/errors"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func CheckFence(loc *storage.Location) {
	openFence, err := storage.CheckFenceSwitch(loc.Imei)
	if err != nil {
		return
	}
	if !openFence {
		return
	}

	fenceInfo, err := storage.GetFence(loc.Imei)
	if err != nil {
		return
	}

	point := storage.Point{
		Lng: float64(loc.Lng) / 1000000,
		Lat: float64(loc.Lat) / 1000000,
	}

	alarm := storage.Alarm{
		Imei:  loc.Imei,
		Time:  time.Unix(loc.Time, 0),
		Lat:   loc.Lat,
		Lng:   loc.Lng,
		Speed: loc.Speed,
	}

	for k, v := range fenceInfo {
		if strings.Contains(v, "Radius") {
			var circle storage.Circle
			json.Unmarshal([]byte(v), &circle)
			in := IsPointInCircle(point, circle)
			if in {
				if circle.FenceState == "2" && (circle.FenceSwitch == "3" || circle.FenceSwitch == "1") {
					alarm.Type = "7"
					alarm.FenceName = circle.FenceName
					storage.InsertAlarm(alarm)
				}
				circle.FenceState = "1"
			} else {
				if circle.FenceState == "1" && (circle.FenceSwitch == "3" || circle.FenceSwitch == "2") {
					alarm.Type = "8"
					alarm.FenceName = circle.FenceName
					storage.InsertAlarm(alarm)
				}
				circle.FenceState = "2"
			}
			fenceValue, _ := json.Marshal(circle)
			storage.SetFence(loc.Imei, k, string(fenceValue))
		} else if strings.Contains(v, "Points") {
			var polygon storage.Polygon
			json.Unmarshal([]byte(v), &polygon)
			in := IsPointInPolygon(point, polygon)
			if in {
				if polygon.FenceState == "2" && (polygon.FenceSwitch == "3" || polygon.FenceSwitch == "1") {
					alarm.Type = "7"
					alarm.FenceName = polygon.FenceName
					storage.InsertAlarm(alarm)
				}
				polygon.FenceState = "1"
			} else {
				if polygon.FenceState == "1" && (polygon.FenceSwitch == "3" || polygon.FenceSwitch == "2") {
					alarm.Type = "8"
					alarm.FenceName = polygon.FenceName
					storage.InsertAlarm(alarm)
				}
				polygon.FenceState = "2"
			}
			fenceValue, _ := json.Marshal(polygon)
			storage.SetFence(loc.Imei, k, string(fenceValue))
		} else if strings.Contains(v, "City") {
			var region storage.Region
			json.Unmarshal([]byte(v), &region)
			in, err := IsPointInRegion(point, region)
			if err != nil {
				log.Warnf("rvsgeo failed")
				continue
			}
			if in {
				if region.FenceState == "2" && (region.FenceSwitch == "3" || region.FenceSwitch == "1") {
					alarm.Type = "7"
					alarm.FenceName = region.FenceName
					storage.InsertAlarm(alarm)
				}
				region.FenceState = "1"
			} else {
				if region.FenceState == "1" && (region.FenceSwitch == "3" || region.FenceSwitch == "2") {
					alarm.Type = "8"
					alarm.FenceName = region.FenceName
					storage.InsertAlarm(alarm)
				}
				region.FenceState = "2"
			}
			fenceValue, _ := json.Marshal(region)
			storage.SetFence(loc.Imei, k, string(fenceValue))
		} else {
			log.Warnf("fence info invalid")
		}
	}
}

func IsPointInPolygon(point storage.Point, polygon storage.Polygon) bool {
	// 初始化交点数量
	intersections := 0
	// 获取多边形的点数
	n := len(polygon.Points)

	// 遍历多边形的边
	for i := 0; i < n; i++ {
		// 获取多边形的当前点和下一个点
		currentVertex := polygon.Points[i]
		nextVertex := polygon.Points[(i+1)%n]

		// 如果目标点与多边形的边平行，则忽略
		if point.Lat < currentVertex.Lat && point.Lat < nextVertex.Lat ||
			point.Lat > currentVertex.Lat && point.Lat > nextVertex.Lat ||
			point.Lng > currentVertex.Lng && point.Lng > nextVertex.Lng {
			continue
		}

		// 计算射线与边的交点
		xIntersection := (point.Lat-currentVertex.Lat)*(nextVertex.Lng-currentVertex.Lng)/(nextVertex.Lat-currentVertex.Lat) + currentVertex.Lng

		// 如果交点在目标点的右边，则增加交点数量
		if xIntersection > point.Lng {
			intersections++
		}
	}

	// 如果交点数量为奇数，则目标点在多边形内部
	return intersections%2 != 0
}

func IsPointInCircle(point storage.Point, circle storage.Circle) bool {
	point1 := geo.NewPoint(point.Lat, point.Lng)
	point2 := geo.NewPoint(circle.Lat, circle.Lng)
	distance := point1.GreatCircleDistance(point2) * 1000
	if distance <= circle.Radius {
		return true
	} else {
		return false
	}
}

func IsPointInRegion(point storage.Point, region storage.Region) (bool, error) {
	type AddressComponent struct {
		City     string `json:"city"`
		Province string `json:"province"`
		Citycode string `json:"citycode"`
		District string `json:"district"`
	}

	type RegeoCode struct {
		AddressComponent AddressComponent `json:"addressComponent"`
	}

	type RvsGeoResult struct {
		Status    string    `json:"status"`
		Info      string    `json:"info"`
		Infocode  string    `json:"infocode"`
		Regeocode RegeoCode `json:"regeocode"`
	}

	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=caeb0f47453c52b1d4136bdaeee1e8dd&output=JSON&location=%v,%v", point.Lng, point.Lat)

	res, err := http.Get(url)
	if err != nil {
		return false, errors.New("http get failed")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, errors.New("parse failed")
	}

	content := string(body)
	fmt.Println(content)

	var result RvsGeoResult
	json.Unmarshal(body, &result)

	if result.Status != "1" || result.Infocode != "10000" {
		return false, errors.New("status error")
	}

	if region.Area != "" {
		if result.Regeocode.AddressComponent.District == region.Area {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		if result.Regeocode.AddressComponent.City == region.City {
			return true, nil
		} else {
			return false, nil
		}
	}
}
