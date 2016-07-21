package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
)

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("data.db", 0600, nil)
	if err != nil {
		return
	}
}

func handlerMux(c *iris.Context) {
	var t msgType
	c.ReadXML(&t)
	switch t.MsgType {
	case MsgTypeText:
		handleText(c)
	case MsgTypeImage:
		handleImage(c)
	case MsgTypeVoice:
		handleVoice(c)
	case MsgTypeVideo:
		handleVideo(c)
	case MsgTypeMusic:
		handleMusic(c)
	case MsgTypeNews:
		handleNews(c)
	default:
		c.Log("不支持该类型，%s.\n", t.MsgType)
		c.WriteString("")
	}
}

func handleText(c *iris.Context) {
	msg := NewText(c.PostBody())
	fmt.Println(msg)

	rmsg := Text{
		ToUserName:   msg.FromUserName,
		FromUserName: developerID,
		CreateTime:   time.Now().Unix(),
	}

	db.Update(func(tx *bolt.Tx) error {
		content := msg.Content
		if strings.HasPrefix(content, "我的学号是") {
			content = content[len(content)-8:]
			b := tx.Bucket([]byte("default"))
			err := b.Put([]byte(msg.FromUserName), []byte(content))
			rmsg.Content = fmt.Sprintf("你的学号是%s\n", content)

			return err
		}

		b := tx.Bucket([]byte("default"))
		data := b.Get([]byte(msg.FromUserName))
		if data == nil {
			rmsg.Content = `请输入"我的学号是00000000"`
			return nil
		}

		if string(data) == "09170515" {
			rmsg.Content = "你是跳跳，一个大美女"
			return nil
		}
		if string(data) == "09170512" {
			rmsg.Content = "你是乐乐，一个大美女"
			return nil
		}
		rmsg.Content = fmt.Sprintf("你的学号是%s，你是%s", data, "我们班的同学")

		return nil
	})

	c.WriteString(rmsg.String())

}

func handleImage(c *iris.Context) {}
func handleVoice(c *iris.Context) {}
func handleVideo(c *iris.Context) {}
func handleMusic(c *iris.Context) {}
func handleNews(c *iris.Context)  {}
