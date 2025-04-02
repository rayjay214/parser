package storage

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	Sms struct {
		Key    string `ini:"key"`
		Secret string `ini:"secret"`
		Sn     string `ini:"sn"`
		Tc     string `ini:"tc"`
	} `ini:"sms"`

	Jt808Server struct {
		Port int `ini:"port"`
	} `ini:"jt808_server"`

	Cassandra struct {
		Host   string `ini:"host"`
		Dbname string `ini:"dbname"`
	} `ini:"cassandra"`

	Mysql struct {
		Host   string `ini:"host"`
		Dbname string `ini:"dbname"`
	} `ini:"mysql"`

	Redis struct {
		Host string `ini:"host"`
	} `ini:"redis"`

	Minio struct {
		Host string `ini:"host"`
	} `ini:"minio"`

	RawLog struct {
		Path string `ini:"path"`
		Open int    `ini:"open"`
	} `ini:"raw_log"`

	Grpc struct {
		Host string `ini:"host"`
	} `ini:"grpc"`
}

var Conf *Config

func LoadConfig(filename string) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return
	}
	Conf = &Config{}
	err = cfg.MapTo(Conf)
	if err != nil {
		return
	}
}
