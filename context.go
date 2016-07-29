package main

import (
	"fmt"
	"net/http"
	"time"
)

//Handler todo
type Handler interface {
	ServeContext(ctx *Context)
}

//HandlerFunc ...
type HandlerFunc func(*Context)

// var _ Handler = HandlerFunc(nil)

// ServeContext todo
func (fn HandlerFunc) ServeContext(ctx *Context) { fn(ctx) }

//Context ...
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	OpenID string
	Type   MsgType

	Message *Text

	index    int
	handlers []Handler
}

// ServeContext todo
func (c *Context) ServeContext(ctx *Context) {
	if ctx != nil {
		c = ctx
	}

	for c.index++; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index].ServeContext(c)
	}
	c.index--
}

func (c *Context) Write(data []byte) (int, error) {
	return c.ResponseWriter.Write(data)
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
