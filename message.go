package main

type MsgType string

const (
	MsgTypeText  MsgType = "text"  // 文本消息
	MsgTypeImage MsgType = "image" // 图片消息
	MsgTypeVoice MsgType = "voice" // 语音消息
	MsgTypeVideo MsgType = "video" // 视频消息
	MsgTypeMusic MsgType = "music" // 音乐消息
	MsgTypeNews  MsgType = "news"  // 图文消息
)

type msgType struct {
	MsgType MsgType `xml:"MsgType"`
}
