package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//Context todo
type Context struct {
	//
	ResponseWriter http.ResponseWriter
	//
	OpenID  string
	Type    string
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
		handlers: make([]Handler, 0, 8),
		buffer:   bytes.NewBuffer(nil),
	}
}

func (c *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
	if !validateURL(signature, timestamp, nonce, cfg.Token) {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

	// 读取数据并Decode
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.ResponseText("数据错误")
		return
	}
	xml.Unmarshal(data, c.Message)

	// 设置Context的值
	c.ResponseWriter = w
	c.OpenID = openid
	c.Type = c.Message.MsgType
	c.index = 0
	c.buffer.Reset()

	// 检查并执行Handlers
	if c.Type != "text" || c.handlers == nil || len(c.handlers) == 0 {
		c.ResponseText("暂不支持此类型信息")
		return
	}
	c.handlers[0].ServeMessage(c)

	// 返回数据
	if c.buffer.Len() <= 0 {
		c.ResponseText("信息格式错误")
		return
	}
	c.ResponseText(c.buffer.String())
}

// Printf todo
func (c *Context) Printf(s string, a ...interface{}) {
	fmt.Fprintf(c.buffer, s, a...)
}

// ResponseText todo
func (c *Context) ResponseText(content string) {
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

// Use todo
func (c *Context) Use(h ...Handler) {
	c.handlers = append(c.handlers, h...)
}

// UseFunc todo
func (c *Context) UseFunc(fns ...func(h *Context)) {
	for _, fn := range fns {
		c.handlers = append(c.handlers, HandlerFunc(fn))
	}

}

// Run todo
func (c *Context) Run() {

	http.ListenAndServe(":80", c)
}

// Next todo
func (c *Context) Next() {
	c.index++
	if c.index >= len(c.handlers) {
		c.index--
		return
	}
	c.handlers[c.index].ServeMessage(c)
}
