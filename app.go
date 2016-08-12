package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

//App todo
type App struct {
	//
	OpenID         string
	ResponseWriter http.ResponseWriter
	//
	commands  map[string]*Command
	outBuffer *bytes.Buffer
	//
	log *logrus.Logger
}

// New todo
func New() *App {
	return &App{
		log:       lg,
		commands:  make(map[string]*Command),
		outBuffer: bytes.NewBuffer(nil),
	}
}

// Run todo
func (c *App) addDefaultCommand() {
	c.Command("我的姓名", func(ctx *Context) {
		if ctx.ArgsLen() != 2 {
			c.Text("我的姓名 XXX")
			return
		}
		if ctx.PersonInfo.Name != "" {
			c.Printf("你的姓名是%s,如错误请联系管理员", ctx.PersonInfo.Name)
			return
		}
		ctx.PersonInfo.OpenID = ctx.Message.FromUserName
		ctx.PersonInfo.Name = ctx.Arg(0)
		ctx.PersonInfo.Put()
		c.Printf("姓名设置成功")
	})

	c.Command("我的学号", func(ctx *Context) {
		if ctx.ArgsLen() != 2 {
			return
		}
		if ctx.PersonInfo.StudentID != "" {
			c.Printf("你的学号是%s,错误请联系管理员", ctx.PersonInfo.StudentID)
			return
		}
		ctx.PersonInfo.OpenID = ctx.Message.FromUserName
		ctx.PersonInfo.StudentID = ctx.Arg(0)
		ctx.PersonInfo.Put()
		c.Printf("学号设置成功")
	})
}

// Run todo
func (c *App) Run() {
	if err := http.ListenAndServe(":80", c); err != nil {
		c.LogWithError(err).Fatal("启动失败")
	}
}

// ServeHTTP todo
func (c *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isValidateRequest(r) {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

	c.ResponseWriter = w
	c.OpenID = r.Form.Get("openid")
	c.outBuffer.Reset()

	ctx, err := NewContext(c, r.Body)
	if err != nil {
		return
	}

	c.addDefaultCommand()
	command, ok := c.commands[ctx.CommandName()]
	if !ok {
		c.Text("不支持")
		return
	}

	switch ctx.CommandName() {
	case "我的姓名", "我的学号":
		command.Action(ctx)
		return
	}

	if ctx.PersonInfo.Name == "" {
		c.Printf(`请输入"我的姓名 XXX"`)
		return
	}
	if ctx.PersonInfo.StudentID == "" {
		c.Printf(`请输入"我的学号 XXXXXXXX"`)
		return
	}

	command.Action(ctx)
	//
	if c.outBuffer.Len() <= 0 {
		c.outBuffer.WriteString("success")
	}
	c.Text(c.outBuffer.String())
}

// Printf todo
func (c *App) Printf(s string, a ...interface{}) {
	fmt.Fprintf(c.outBuffer, s, a...)
}

// Text todo
func (c *App) Text(content string) {
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

// LogWithError todo
func (c *App) LogWithError(err error) *logrus.Entry {
	return c.log.WithError(err)
}

// LogWithFiled todo
func (c *App) LogWithFiled(key string, value interface{}) *logrus.Entry {
	return c.log.WithField(key, value)
}

// Command todo
func (c *App) Command(name string, fn func(*Context)) {
	c.commands[name] = &Command{Action: fn}
}
