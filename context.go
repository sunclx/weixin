package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

//Context todo
type Context struct {
	//
	ResponseWriter http.ResponseWriter

	//
	OpenID     string
	Message    *Text
	Command    *Command
	PersonInfo *PersonInfo
	//
	handler   Handler
	inBuffer  *bytes.Buffer
	outBuffer *bytes.Buffer

	err error
	log *logrus.Logger
}

// New todo
func New() *Context {
	return &Context{
		log:     lg,
		Message: new(Text),

		inBuffer:  bytes.NewBuffer(nil),
		outBuffer: bytes.NewBuffer(nil),
	}
}

// Run todo
func (c *Context) Run() {
	if err := http.ListenAndServe(":80", c); err != nil {
		c.LogWithError(err).Fatal("启动失败")
	}
}

// ServeHTTP todo
func (c *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查域名及请求方法
	if hostname := r.Host; r.Method != "POST" || hostname != "weixin.chenlixin.net" || r.URL.Path != "/" {
		fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
		w.Write([]byte("404"))
		return
	}

	// 检验请求参数
	r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	openid := r.Form.Get("openid")
	if !validateURL(signature, timestamp, nonce, cfg.Token) {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

	// 设置Context的值
	c.ResponseWriter = w
	c.OpenID = openid
	c.outBuffer.Reset()

	c.inBuffer.Reset()
	c.inBuffer.ReadFrom(r.Body)
	xml.Unmarshal(c.inBuffer.Bytes(), c.Message)
	if c.Message.MsgType != "text" {
		c.Text("暂不支持此类型信息")
		return
	}
	if c.Message.FromUserName != c.OpenID {
		c.Text("微信服务器错误，请稍后再试")
		return
	}

	// Default
	command := NewCommand(c.Message.Content)
	err := c.PersonInfo.Get(c.OpenID)
	if err != nil {
		c.LogWithError(err).Infof("获取openid: %s个人信息错误", c.OpenID)
		return
	}

	switch command.Name {
	case "我的姓名":
		if len(command.Arguments) != 2 {
			return
		}
		if c.PersonInfo.Name != "" {
			c.Printf("你的姓名是%s,如错误请联系管理员", c.PersonInfo.Name)
			return
		}
		c.PersonInfo.OpenID = c.OpenID
		c.PersonInfo.Name = command.Argument(1)
		c.PersonInfo.Put()
		c.Printf("姓名设置成功")
		return
	case "我的学号":
		if len(command.Arguments) != 2 {
			return
		}
		if c.PersonInfo.StudentID != "" {
			c.Printf("你的学号是%s,错误请联系管理员", c.PersonInfo.StudentID)
			return
		}
		c.PersonInfo.OpenID = c.OpenID
		c.PersonInfo.StudentID = command.Argument(1)
		c.PersonInfo.Put()
		c.Printf("学号设置成功")
		return
	}

	if c.PersonInfo.Name == "" {
		c.Printf(`请输入"我的姓名 XXX"`)
		return
	}
	if c.PersonInfo.StudentID == "" {
		c.Printf(`请输入"我的学号 XXXXXXXX"`)
		return
	}

	// 调用handlers
	c.handler.ServeMessage(c)
	if c.outBuffer.Len() <= 0 {
		c.outBuffer.WriteString("success")
	}
	c.Text(c.outBuffer.String())
}

// Printf todo
func (c *Context) Printf(s string, a ...interface{}) {
	fmt.Fprintf(c.outBuffer, s, a...)
}

// Text todo
func (c *Context) Text(content string) {
	fmt.Fprintf(c.ResponseWriter, `
	<xml>
	<ToUserName><![CDATA[%s]]></ToUserName>
	<FromUserName><![CDATA[%s]]></FromUserName>
	<CreateTime>%d</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	<Content><![CDATA[%s]]></Content>
	</xml>`,
		c.OpenID, cfg.DeveloperID, time.Now().Unix(), content)
}

// Handler todo
func (c *Context) Handler(h Handler) {
	c.handler = h
}

// LogWithError todo
func (c *Context) LogWithError(err error) *logrus.Entry {
	return c.log.WithError(err)
}

// LogWithFiled todo
func (c *Context) LogWithFiled(key string, value interface{}) *logrus.Entry {
	return c.log.WithField(key, value)
}
