package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

//Context todo
type Context struct {
	//
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	//
	appID     string
	token     string
	secruteID string
	//
	OpenID  string
	Type    MsgType
	Message *Text
	//
	index    int
	handlers []Handler
	//
	buffer *bytes.Buffer
}

// New todo
func New() *Context {
	return &Context{
		//
		appID:     cfg.AppID,
		token:     cfg.Token,
		secruteID: cfg.SecruteID,
		//
		handlers: make([]Handler, 0, 8),
		buffer:   bytes.NewBuffer(nil),
	}
}

func (s *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查域名及请求方法
	if hostname := r.Host; r.Method != "POST" || hostname != "weixin.chenlixin.net" || r.URL.Path != "/" {
		fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
		w.Write([]byte("404"))
		return
	}

	// 检验请求参数
	r.ParseForm()
	queryParams := r.Form
	signature := queryParams.Get("signature")
	timestamp := queryParams.Get("timestamp")
	nonce := queryParams.Get("nonce")
	openid := queryParams.Get("openid")
	if !validateURL(signature, timestamp, nonce, s.token) {
		w.Write([]byte("404"))
		return
	}

	buffer := bytes.NewBuffer(nil)
	buffer.Reset()
	buffer.ReadFrom(r.Body)

	var t Text
	xml.Unmarshal(buffer.Bytes(), &t)

	ctx := &Context{
		ResponseWriter: w,
		Request:        r,

		appID:     cfg.AppID,
		token:     cfg.Token,
		secruteID: cfg.SecruteID,

		OpenID: openid,
		Type:   t.MsgType,

		Message: &t,
	}

	if ctx.handlers == nil || len(ctx.handlers) == 0 {
		return
	}
	ctx.handlers[0].ServeMessage(ctx)

	if ctx.Message.MsgType != MsgTypeText {
		ctx.ResponseText("暂不支持此类型信息")
		return
	}

	if ctx.buffer.Len() <= 0 {
		ctx.ResponseText("你的信息格式错误")
		return
	}
	ctx.ResponseText(ctx.buffer.String())

}

func (c *Context) Printf(s string, a ...interface{}) {
	fmt.Fprintf(c.buffer, s, a...)
}

// ResponseText todo
func (c *Context) ResponseText(content string) {
	c.Printf(`
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`,
		c.OpenID, cfg.DeveloperID, time.Now().Unix(), content)

}

func (c *Context) Use(h ...Handler) {
	c.handlers = append(c.handlers, h...)
}

func (c *Context) UseFunc(fns ...func(h *Context)) {
	for _, fn := range fns {
		c.handlers = append(c.handlers, HandlerFunc(fn))
	}

}

func (c *Context) Run() {

	http.ListenAndServe(":80", c)
}

func (c *Context) Next() {
	c.index++
	if c.index >= len(c.handlers) {
		c.index--
		return
	}
	c.handlers[c.index].ServeMessage(c)
}
