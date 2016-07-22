package main

import (
	"fmt"
	"strings"

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
		rcontent = handlePhone(content)
	default:
		rcontent = "你的消息格式暂不支持"
	}

	c.WriteString(RText(openid, rcontent).String())

}

func handleImage(c *iris.Context) {}
func handleVoice(c *iris.Context) {}
func handleVideo(c *iris.Context) {}
func handleMusic(c *iris.Context) {}
func handleNews(c *iris.Context)  {}
