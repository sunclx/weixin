package main

import "github.com/boltdb/bolt"

var db *bolt.DB
var cfg *config

func main() {
	// 初始化配置
	initConfig()

	//设置数据库
	var err error
	db, err = bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		buckets := []string{"default", "Contact", "Phone", "Person", "NameID"}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})

	//启动服务
	s := New()
	s.UseFunc(logHandler)
	s.UseFunc(testHandler)
	s.Run(":80")
}

type config struct {
	DeveloperID string
	AppID       string
	Token       string
	SecruteID   string
	DBPath      string
}

func initConfig() {
	cfg = &config{
		DeveloperID: "gh_3fb3b0b8f2fa",
		AppID:       "",
		Token:       "njmu0917",
		SecruteID:   "",
		DBPath:      "/root/data.db",
	}
}
