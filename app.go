package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

//Cli 是基本类型
type Cli struct {
	handlers []http.Handler
	commands map[string]*Command
}

// New 返回一个新 *Cli
func New() *Cli {
	return &Cli{
		commands: make(map[string]*Command),
	}
}

// Run 运行Cli
func (c *Cli) Run() {
	if err := http.ListenAndServe(":80", c); err != nil {
		log.WithError(err).Fatal("启动失败")
	}
}

// ServeHTTP 实现了htto.Handler
func (c *Cli) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, handler := range c.handlers {
		handler.ServeHTTP(w, r)
		return
	}
	// 验证请求
	if !isValidateRequest(r) {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

	// 返回数据
	out := bufferPool.Get().(*bytes.Buffer)
	out.Reset()
	defer func() {
		if out.Len() <= 0 {
			out.WriteString("failed")
		}

		fmt.Fprintf(w, `
	<xml>
	<ToUserName><![CDATA[%s]]></ToUserName>
	<FromUserName><![CDATA[%s]]></FromUserName>
	<CreateTime>%d</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	<Content><![CDATA[%s]]></Content>
	</xml>`,
			r.Form.Get("openid"), cfg.DeveloperID, time.Now().Unix(), out.String())
	}()

	ctx, err := NewContext(out, r.Body)
	if err != nil {
		return
	}
	log.WithField("openid", ctx.Message.FromUserName).Infof("%#v\n", ctx)

	// 执行Command
	switch command, ok := c.commands[ctx.CommandName()]; {
	case ctx.CommandName() == "我的姓名", ctx.CommandName() == "我的学号":
		command.Run(ctx)
	case ctx.User.Name == "":
		ctx.Print(`请输入"我的姓名 XXX"`)
	case ctx.User.StudentID == "":
		ctx.Print(`请输入"我的学号 XXXXXXXX"`)
	case ok:
		command.Run(ctx)
	}
}

// Use 添加一个新命令
func (c *Cli) Use(h http.Handler) {
	c.handlers = append(c.handlers, h)
}

// UseFunc 添加一个新命令
func (c *Cli) UseFunc(fn func(w http.ResponseWriter, r *http.Request)) {
	c.handlers = append(c.handlers, http.HandlerFunc(fn))
}

// Command 添加一个新命令
func (c *Cli) Command(name string, fn func(*Context)) {
	c.commands[name] = NewCommand(fn)
}
