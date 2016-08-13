package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Context 是Command的上下文
type Context struct {
	w io.Writer

	commandName      string
	commandArguments []string

	Message Text
	User    User
}

// NewContext 读取io.Reader返回*Context，nil；若数据错误返回nil,error
func NewContext(w io.Writer, r io.Reader) (*Context, error) {
	var ctx Context
	ctx.w = w

	// 解析message
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.ReadFrom(r)
	if err := xml.Unmarshal(buffer.Bytes(), &ctx.Message); err != nil {
		return nil, err
	}
	if ctx.Message.MsgType != "text" {
		return nil, errors.New("暂不支持此类型信息")
	}

	// 获取User
	if err := ctx.User.Get(ctx.Message.FromUserName); err != nil {
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

// CommandName todo
func (c *Context) CommandName() string {
	return c.commandName
}

// NArg todo
func (c *Context) NArg() int {
	return len(c.commandArguments)
}

// Arg todo
func (c *Context) Arg(index int) string {
	if index-1 >= c.NArg() || index-1 < 0 {
		return ""
	}
	return c.commandArguments[index-1]
}

// Print todo
func (c *Context) Print(a ...interface{}) {
	fmt.Fprint(c.w, a...)
}

// Printf todo
func (c *Context) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.w, format, a...)
}
