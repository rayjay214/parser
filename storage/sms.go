package storage

import (
	"github.com/aliyun-sdk/sms-go"
	"gopkg.in/ini.v1"
	"log"
	"strings"
)

var err error
var client *sms.Client
var ak, sk, sn, tc string

func init() {
	cfg, err := ini.Load("config.ini")
	section := cfg.Section("general")

	ak = section.Key("key").Value()
	sk = section.Key("secret").Value()
	sn = section.Key("sn").Value()
	tc = section.Key("tc").Value()
	//sn = "千讯互联"
	//tc = "SMS_262415663"
	client, err = sms.New(ak, sk, sms.SignName(sn), sms.Template(tc))
	if err != nil {
		log.Fatalln(err)
	}
}

func Send(phone string, imei string, alarm string) error {
	alarmPrint := strings.Replace(alarm, "报警", "", -1)
	err = client.Send(
		sms.Mobile(phone),
		sms.Parameter(map[string]string{
			"imei":  imei,
			"alarm": alarmPrint,
		}),
	)
	if err != nil {
		log.Println(err)
	}
	return err
}
