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

// HandlerFunc todo
type HandlerFunc func(c *Context)

// ServeMessage todo
func (fn HandlerFunc) ServeMessage(c *Context) {
	fn(c)
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
func (p *PersonInfo) Put(openid string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("PersonInfo"))
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return bx.Put([]byte(openid), data)
	})
}

type contactHandler struct {
	PersonInfo PersonInfo
}

// ServeMessage todo
func (c *contactHandler) ServeMessage(ctx *Context) {
	parts := strings.Fields(ctx.Message.Content)
	err := c.PersonInfo.Get(ctx.OpenID)
	// if true {
	// 	ctx.Printf("success")
	// 	return
	// }
	if err != nil {
		ctx.Printf("服务器错误")
		ctx.WithError(err).Errorln("获取个人信息错误")
		return
	}

	switch parts[0] {
	case "我的姓名":
		if len(parts) != 2 {
			ctx.Printf("参数错误")
			ctx.Infoln(parts)
			return
		}

		c.PersonInfo.Name = parts[1]
		c.PersonInfo.Put(ctx.OpenID)
		db.Update(func(tx *bolt.Tx) error {
			bx, _ := tx.CreateBucketIfNotExists([]byte("NameOpenID"))
			return bx.Put([]byte(c.PersonInfo.Name), []byte(ctx.OpenID))
		})
		ctx.Printf("设置成功")
		return
	case "我的学号":
		if len(parts) != 2 {
			ctx.Printf("参数错误")
			ctx.Infoln(parts)
			return
		}

		c.PersonInfo.StudentID = parts[1]
		c.PersonInfo.Put(ctx.OpenID)
		ctx.Printf("设置成功")
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

	switch parts[0] {
	case "手机":
		if len(parts) != 2 {
			ctx.Printf(`
			信息格式错误
			输入"我的姓名 XXX"设置姓名
			输入"我的学号 XXXXXXXX"设置学号
			输入"我的手机 XXX"设置手机
			输入"手机 姓名"查询手机号码
			`)
			ctx.Infoln(parts)
			return
		}
		name := parts[1]
		openid := ""
		db.Update(func(tx *bolt.Tx) error {
			bx, _ := tx.CreateBucketIfNotExists([]byte("NameOpenID"))
			openid = string(bx.Get([]byte(name)))
			return nil
		})
		err := c.PersonInfo.Get(openid)
		if err != nil {
			ctx.Printf(`服务器错误`)
			ctx.WithError(err).Errorln("获取个人信息错误")
			return
		}
		if c.PersonInfo.PhoneNumber == "" {
			ctx.Printf("没有%s的号码", name)
			return
		}

		ctx.Printf("%s %s", c.PersonInfo.Name, c.PersonInfo.PhoneNumber)
	case "我的手机":
		if len(parts) != 2 {
			ctx.Printf(`
			信息格式错误
			输入"我的姓名 XXX"设置姓名
			输入"我的学号 XXXXXXXX"设置学号
			输入"我的手机 XXX"设置手机
			输入"手机 姓名"查询手机号码
			`)
			ctx.Infoln(parts)
			return
		}

		c.PersonInfo.PhoneNumber = parts[1]
		c.PersonInfo.Put(ctx.OpenID)

		ctx.Printf("设置成功")
	case "绑定姓名":
		if len(parts) != 2 {
			ctx.Printf("参数错误")
			ctx.Infoln(parts)
			return
		}

	default:
		ctx.Next()
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

// type NameID struct {
// 	Name      string
// 	StudentID string
// }

// func (n *NameID) Get(openid string) error {
// 	bt, err := bx.New([]byte("NameID"))
// 	if err != nil {
// 		return err
// 	}

// 	data, err := bt.Get([]byte(openid))
// 	if data == nil || err != nil {
// 		return errors.New("Not Exist")
// 	}

// 	parts := bytes.Split(data, []byte("&&"))
// 	if len(parts) != 2 {
// 		return errors.New("wrong data")
// 	}

// 	n.Name = string(parts[0])
// 	n.StudentID = string(parts[1])
// 	return nil
// }
// func (n *NameID) Put(openid string) error {
// 	bt, err := bx.New([]byte("NameID"))
// 	if err != nil {
// 		return err
// 	}

// 	return bt.Put([]byte(openid), []byte(n.Name+"&&"+n.StudentID))
// }

// type openidHandler struct {
// 	NameID
// }

// func (o *openidHandler) ServeMessage(ctx *Context) {
// 	content := ctx.Message.Content
// 	parts := strings.Fields(content)
// 	switch parts[0] {
// 	case "我的学号":
// 		if len(parts) != 2 {
// 			ctx.Printf("参数错误")
// 			return
// 		}
// 		o.Name = parts[1]
// 		o.StudentID = parts[1]
// 		o.Put(ctx.OpenID)
// 		ctx.Printf("学号绑定成功")
// 		return
// 	}

// 	err := o.Get(ctx.OpenID)
// 	if err != nil {
// 		ctx.Printf(`请输入“我的学号 00000000”`)
// 		return
// 	}

// 	ctx.Next()
// }
