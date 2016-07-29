package main

import (
	"bytes"
	"fmt"
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

type Message struct {
	msg    *Text
	buffer *bytes.Buffer

	index           int
	messageHandlers []MessageHandler
}

// func (m *Message) Write(data []byte) (int, error) {
// 	if m.buffer == nil {
// 		m.buffer = bytes.NewBuffer(nil)
// 	}

// 	return m.buffer.Write(data)
// }

func (m *Message) Printf(s string, a ...interface{}) {
	fmt.Fprintf(m.buffer, s, a...)
}
func (m *Message) Use(h ...MessageHandler) *Message {
	if m.messageHandlers == nil {
		m.messageHandlers = make([]MessageHandler, 0, 8)
	}
	m.messageHandlers = append(m.messageHandlers, h...)

	return m
}

func (m *Message) UseFunc(fns ...func(msg *Message)) *Message {
	if m.messageHandlers == nil {
		m.messageHandlers = make([]MessageHandler, 0, 8)
	}
	for _, fn := range fns {
		m.messageHandlers = append(m.messageHandlers, MessageHandlerFunc(fn))
	}

	return m
}

func (m *Message) Begin() {
	if m.messageHandlers == nil || len(m.messageHandlers) == 0 {
		return
	}
	m.messageHandlers[0].ServeMessage(m)

}

func (m *Message) Next() {
	m.index++

	if m.index >= len(m.messageHandlers) {
		m.index--
		return
	}
	m.messageHandlers[m.index].ServeMessage(m)
}

type MessageHandler interface {
	ServeMessage(msg *Message)
}

type MessageHandlerFunc func(msg *Message)

func (m MessageHandlerFunc) ServeMessage(msg *Message) {
	m(msg)
}
