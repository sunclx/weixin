package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

// Context todo
type Context struct {
	app *App

	commandName      string
	commandArguments []string

	Message    Text
	PersonInfo PersonInfo
}

func NewContext(app *App, r io.Reader) (*Context, error) {
	var ctx Context
	ctx.app = app

	// 解析message
	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(r)
	if err := xml.Unmarshal(buffer.Bytes(), &ctx.Message); err != nil {
		return nil, err
	}
	if ctx.Message.MsgType != "text" {
		return nil, errors.New("暂不支持此类型信息")
	}

	// 获取PersonInfo
	if err := ctx.PersonInfo.Get(ctx.Message.FromUserName); err != nil {
		return nil, err
	}

	// 获取Command
	ss := strings.Fields(ctx.Message.Content)
	ctx.commandName = ss[0]
	if len(ss) > 1 {
		ctx.commandArguments = ss[1:]
	} else {
		ctx.commandArguments = []string{}
	}

	return &ctx, nil
}
func (c *Context) CommandName() string {
	return c.commandName
}

func (c *Context) ArgsLen() int {
	return len(c.commandArguments)
}

func (c *Context) Arg(index int) string {
	if index >= c.ArgsLen() || index < 0 {
		return ""
	}
	return c.commandArguments[index]
}

// Command todo
type Command struct {
	Action func(*Context)
}
