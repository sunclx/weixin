package main

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/naoina/toml"
)

var (
	configFile = "/root/weixin/config.toml"
	cfg        config
	//	bx  *buckets.DB
	db  *bolt.DB
	log = logrus.New()
)

type config struct {
	DeveloperID string
	AppID       string
	Token       string
	SecruteID   string
	DBPath      string
}

func init() {
	var err error

	//初始化logger
	log.Out, err = os.OpenFile("/root/weixin/weixin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	fatalError(err)

	//初始化配置
	buf, err := ioutil.ReadFile(configFile)
	fatalError(err)

	err = toml.Unmarshal(buf, &cfg)
	fatalError(err)

	//设置数据库
	db, err = bolt.Open(cfg.DBPath, 0600, nil)
	fatalError(err)

}

// Text todo
type Text struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        string `xml:"MsgId"`
}
