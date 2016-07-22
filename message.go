package main

import (
	"encoding/xml"
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

const developerID string = "gh_3fb3b0b8f2fa"

type msgType struct {
	MsgType MsgType `xml:"MsgType"`
}

type Text struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        int64  `xml:"MsgId"`
}

func ParseText(data []byte) Text {
	var t Text
	xml.Unmarshal(data, &t)
	return t

}
func RText(openid, content string) Text {
	return Text{
		ToUserName:   openid,
		FromUserName: developerID,
		MsgType:      "text",
		Content:      content,
	}

}
func (t Text) Marshal() []byte {
	data, _ := xml.Marshal(t)
	return data
}
func (t Text) String() string {
	var templateText = `
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`
	return fmt.Sprintf(templateText, t.ToUserName, t.FromUserName, time.Now().Unix(), t.Content)

}
