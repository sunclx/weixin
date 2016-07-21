package main

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Text struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        int64  `xml:"MsgId"`
}

func unmarshalMsg(data []byte) (msg Text, err error) {
	xml.Unmarshal(data, &msg)
	return msg, nil
}

var templateText = `
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[gh_3fb3b0b8f2fa]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`

func marshaMsg(customerID string, content string) (msg string) {
	return fmt.Sprintf(templateText, customerID, time.Now().Unix(), content)
}
