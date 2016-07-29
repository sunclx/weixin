package main

import "github.com/boltdb/bolt"

var db *bolt.DB
var cfg *config

func init() {
	//设置数据库
	db, _ = bolt.Open(cfg.DBPath, 0600, nil)

}

type config struct {
	DeveloperID string
	AppID       string
	Token       string
	SecruteID   string
	DBPath      string
}

func init() {
	//设置数据库
	cfg = &config{
		DeveloperID: "gh_3fb3b0b8f2fa",
		AppID:       "",
		Token:       "njmu0917",
		SecruteID:   "",
		DBPath:      "/root/data.db",
	}
}
