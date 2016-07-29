package main

import (
	"bytes"
	"fmt"
)

func logHandler(c *Context) {
	r := c.Request
	fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
}

func messageHandler(c *Context) {
	if c.Message.MsgType != MsgTypeText {
		c.ResponseText("暂不支持此类型信息")
		return
	}

	msg := &Message{
		msg:             c.Message,
		buffer:          bytes.NewBuffer(nil),
		messageHandlers: make([]MessageHandler, 0, 8),
	}
	msg.UseFunc(handlePhone)
	msg.UseFunc(handleBindPhone)
	msg.Begin()

	if msg.buffer.Len() <= 0 {
		c.ResponseText("你的信息格式错误")
		return
	}
	c.ResponseText(msg.buffer.String())
}
