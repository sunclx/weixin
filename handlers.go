package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
)

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
	msg := ParseText(c.PostBody())
	fmt.Println(msg)

	openid := msg.FromUserName
	content := msg.Content

	var rcontent string
	switch {
	case strings.HasPrefix(content, PrefixPhone):
		rcontent = handlePhone(msg)
	case strings.HasPrefix(content, PrefixBindPhone):
		rcontent = handleBindPhone(msg)
	default:
		rcontent = "你的消息格式暂不支持"
	}

	c.WriteString(RText(openid, rcontent).String())
}

type person struct {
	StudentID  string
	Name       string
	Birthday   time.Time
	BirthPlace string
	Location   string
}

func personByByte(data []byte) *person {
	var p person
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil
	}
	return &p
}

func (p *person) JSON() string {
	s, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(s)
}
func (p *person) Get() {
	if p.StudentID == "" {
		p = nil
		return
	}

	err := db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("persons"))
		bp := bx.Get([]byte(p.StudentID))
		p = personByByte(bp)
		return nil
	})
	if err != nil {
		p = nil
		return
	}
}
func (p *person) Put() {
	if p.StudentID == "" {
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("persons"))
		return bx.Put([]byte(p.StudentID), []byte(p.JSON()))
	})

}

func handleImage(c *iris.Context) {}
func handleVoice(c *iris.Context) {}
func handleVideo(c *iris.Context) {}
func handleMusic(c *iris.Context) {}
func handleNews(c *iris.Context)  {}
