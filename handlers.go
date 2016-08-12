package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// PersonInfo todo
type PersonInfo struct {
	OpenID string

	StudentID string
	Name      string

	PhoneNumber string
}

//Get todo
func (p *PersonInfo) Get(openid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("PersonInfo"))
		data := bx.Get([]byte(openid))
		if data == nil || string(data) == "" {
			return nil
		}
		return json.Unmarshal(data, p)
	})

}

// Put todo
func (p *PersonInfo) Put() error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("PersonInfo"))
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return bx.Put([]byte(p.OpenID), data)
	})
}

// type contactHandler struct {
// 	PersonInfo *PersonInfo
// }

// }

// type NameOpenID map[string]string

// func (n *NameOpenID) ServeMessage(c *App) {
// 	db.Update(func(tx *bolt.Tx) error {
// 		bx := tx.Bucket([]byte("NameOpenID"))
// 		items := NameOpenID{
// 			"": "",
// 		}

// 		var err error
// 		for k, v := range items {
// 			err = bx.Put([]byte(k), []byte(v))
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

// type NameStudentID map[string]string

// func (n *NameStudentID) ServeMessage(c *App) {
// 	db.Update(func(tx *bolt.Tx) error {
// 		bx := tx.Bucket([]byte("NameOpenID"))
// 		items := NameOpenID{
// 			"": "",
// 		}

// 		var err error
// 		for k, v := range items {
// 			err = bx.Put([]byte(k), []byte(v))
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

// 信息格式错误
// 输入"我的姓名 XXX"设置姓名
// 输入"我的学号 XXXXXXXX"设置学号
// 输入"我的手机 XXX"设置手机
// 输入"手机 姓名"查询手机号码
// 			`)
