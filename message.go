package main

//MsgType todo
type MsgType string

const (
	//MsgTypeText todo
	MsgTypeText MsgType = "text" // 文本消息
//	MsgTypeImage MsgType = "image" // 图片消息
//	MsgTypeVoice MsgType = "voice" // 语音消息
//	MsgTypeVideo MsgType = "video" // 视频消息
//	MsgTypeMusic MsgType = "music" // 音乐消息
//	MsgTypeNews  MsgType = "news"  // 图文消息
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
