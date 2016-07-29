package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

//Context todo
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	appID     string
	token     string
	secruteID string

	OpenID string
	Type   MsgType

	Message *Text
}

// New todo
func New() *Context {
	return &Context{
		appID:     cfg.AppID,
		token:     cfg.Token,
		secruteID: cfg.SecruteID,
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

	if ctx.Message.MsgType != MsgTypeText {
		ctx.ResponseText("暂不支持此类型信息")
		return
	}

	msg := &Message{
		msg:             ctx.Message,
		buffer:          bytes.NewBuffer(nil),
		messageHandlers: make([]MessageHandler, 0, 8),
	}
	msg.UseFunc(handlePhone)
	msg.UseFunc(handleBindPhone)
	msg.Begin()

	if msg.buffer.Len() <= 0 {
		ctx.ResponseText("你的信息格式错误")
		return
	}
	ctx.ResponseText(msg.buffer.String())

}

func (c *Context) Run() {
	http.ListenAndServe(":80", c)
}

// Printf todo
func (c *Context) Printf(s string, a ...interface{}) {
	fmt.Fprintf(c.ResponseWriter, s, a...)
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

//验证函数
func validateURL(signature, timestamp, nonce, token string) bool {
	//排序参数并合并
	ss := []string{token, timestamp, nonce}
	sort.Strings(ss)
	s := strings.Join(ss, "")

	//计算sha1的值
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	//比较计算的signature与获取值比较
	if signatureHex := fmt.Sprintf("%x", bs); signatureHex != signature {
		return false
	}
	return true
}
