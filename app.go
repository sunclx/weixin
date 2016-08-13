package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

//App todo
type App struct {
	commands map[string]*Command
}

// New todo
func New() *App {
	return &App{
		commands: make(map[string]*Command),
	}
}

// Run todo
func (c *App) Run() {
	if err := http.ListenAndServe(":80", c); err != nil {
		log.WithError(err).Fatal("启动失败")
	}
}

// ServeHTTP todo
func (c *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isValidateRequest(r) {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

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
	log.WithField("openid", ctx.Message.FromUserName).Infof("%#v\n", ctx.Message)

	command, ok := c.commands[ctx.CommandName()]
	if !ok {
		out.WriteString("不支持")
		return
	}

	switch ctx.CommandName() {
	case "我的姓名", "我的学号":
		command.Action(ctx)
		return
	}

	if ctx.PersonInfo.Name == "" {
		out.WriteString(`请输入"我的姓名 XXX"`)
		return
	}
	if ctx.PersonInfo.StudentID == "" {
		out.WriteString(`请输入"我的学号 XXXXXXXX"`)
		return
	}

	command.Action(ctx)
}

// Command todo
func (c *App) Command(name string, fn func(*Context)) {
	c.commands[name] = &Command{Action: fn}
}
