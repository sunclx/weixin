package main

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/naoina/toml"
)

type config struct {
	DeveloperID string
	AppID       string
	Token       string
	SecruteID   string
	DBPath      string
}

var (
	configFile = "/root/weixin/config.toml"
	cfg        config
	db         *bolt.DB
	log        = logrus.New()
)

func fatalError(err error) {
	if err != nil {
		log.WithError(err).Fatal("程序错误")
	}
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
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("NameOpenID"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("User"))
		return err
	})

}
