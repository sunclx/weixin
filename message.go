package main

import (
	"fmt"
	"time"
)

type MsgType string

const (
	MsgTypeText  MsgType = "text"  // 文本消息
	MsgTypeImage MsgType = "image" // 图片消息
	MsgTypeVoice MsgType = "voice" // 语音消息
	MsgTypeVideo MsgType = "video" // 视频消息
	MsgTypeMusic MsgType = "music" // 音乐消息
	MsgTypeNews  MsgType = "news"  // 图文消息
)

// Text todo
type Text struct {
	ToUserName   string  `xml:"ToUserName"`
	FromUserName string  `xml:"FromUserName"`
	CreateTime   int64   `xml:"CreateTime"`
	MsgType      MsgType `xml:"MsgType"`
	Content      string  `xml:"Content"`
	MsgID        string  `xml:"MsgId"`
}

func (t Text) String() string {
	return fmt.Sprintf(`
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
<MsgId><![CDATA[%s]]></MsgId>
</xml>`,
		t.ToUserName, t.FromUserName, time.Now().Unix(), t.Content, t.MsgID)

}

// ResponseText todo
func ResponseText(toUserName, fromUserName, content string) string {
	return fmt.Sprintf(`
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`,
		toUserName, fromUserName, time.Now().Unix(), content)

}
