package main

import (
	"encoding/json"
	"strings"

	"github.com/boltdb/bolt"
)

// Handler todo
type Handler interface {
	ServeMessage(c *Context)
}

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

type defaultHandler struct {
	PersonInfo *PersonInfo
}

func (c *defaultHandler) ServeMessage(ctx *Context) {
	paramters := strings.Fields(ctx.Message.Content)
	command := paramters[0]
	err := c.PersonInfo.Get(ctx.OpenID)
	if err != nil {
		ctx.LogWithError(err).Infof("获取openid: %s个人信息错误", ctx.OpenID)
		return
	}

	switch command {
	case "我的姓名":
		if len(paramters) != 2 {
			return
		}
		if c.PersonInfo.Name != "" {
			ctx.Printf("你的姓名是%s,如错误请联系管理员", c.PersonInfo.Name)
			return
		}
		c.PersonInfo.OpenID = ctx.OpenID
		c.PersonInfo.Name = paramters[1]
		c.PersonInfo.Put()
		ctx.Printf("姓名设置成功")
		return
	case "我的学号":
		if len(paramters) != 2 {
			return
		}
		if c.PersonInfo.StudentID != "" {
			ctx.Printf("你的学号是%s,错误请联系管理员", c.PersonInfo.StudentID)
			return
		}
		c.PersonInfo.OpenID = ctx.OpenID
		c.PersonInfo.StudentID = paramters[1]
		c.PersonInfo.Put()
		ctx.Printf("学号设置成功")
		return
	}

	if c.PersonInfo.Name == "" {
		ctx.Printf(`请输入"我的姓名 XXX"`)
		return
	}
	if c.PersonInfo.StudentID == "" {
		ctx.Printf(`请输入"我的学号 XXXXXXXX"`)
		return
	}
}

type contactHandler struct {
	PersonInfo *PersonInfo
}

// ServeMessage todo
func (c *contactHandler) ServeMessage(ctx *Context) {
	paramters := strings.Fields(ctx.Message.Content)
	switch paramters[0] {
	case "手机":
		if len(paramters) != 2 {
			ctx.Printf(`
信息格式错误
输入"我的姓名 XXX"设置姓名
输入"我的学号 XXXXXXXX"设置学号
输入"我的手机 XXX"设置手机
输入"手机 姓名"查询手机号码
			`)
			return
		}
		name := paramters[1]
		openid := ""
		db.Update(func(tx *bolt.Tx) error {
			bx, _ := tx.CreateBucketIfNotExists([]byte("NameOpenID"))
			openid = string(bx.Get([]byte(name)))
			return nil
		})
		p := PersonInfo{}
		err := p.Get(openid)
		if err != nil {
			ctx.Printf(`服务器错误`)
			ctx.LogWithError(err).Errorln("获取个人信息错误")
			return
		}
		if p.PhoneNumber == "" {
			ctx.Printf("没有%s的号码", name)
			return
		}

		ctx.Printf("%s %s", p.Name, p.PhoneNumber)
	case "我的手机":
		if len(paramters) != 2 {
			ctx.Printf(`
信息格式错误
输入"我的姓名 XXX"设置姓名
输入"我的学号 XXXXXXXX"设置学号
输入"我的手机 XXX"设置手机
输入"手机 姓名"查询手机号码
			`)
			return
		}
		c.PersonInfo.OpenID = ctx.OpenID
		c.PersonInfo.PhoneNumber = paramters[1]
		c.PersonInfo.Put()

		ctx.Printf("设置成功")

		// default:
		// 	ctx.Next()
	}

}

type NameOpenID map[string]string

func (n *NameOpenID) ServeMessage(c *Context) {
	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("NameOpenID"))
		items := NameOpenID{
			"": "",
		}

		var err error
		for k, v := range items {
			err = bx.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type NameStudentID map[string]string

func (n *NameStudentID) ServeMessage(c *Context) {
	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("NameOpenID"))
		items := NameOpenID{
			"": "",
		}

		var err error
		for k, v := range items {
			err = bx.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// 信息格式错误
// 输入"我的姓名 XXX"设置姓名
// 输入"我的学号 XXXXXXXX"设置学号
// 输入"我的手机 XXX"设置手机
// 输入"手机 姓名"查询手机号码
// 			`)
