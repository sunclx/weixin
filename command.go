package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/boltdb/bolt"
)

// Context todo
type Context struct {
	w io.Writer

	commandName      string
	commandArguments []string

	Message    Text
	PersonInfo PersonInfo
}

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

func (c *Context) Print(a ...interface{}) {
	fmt.Fprint(c.w, a...)
}

func (c *Context) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.w, format, a...)
}

// Command todo
type Command struct {
	Action func(*Context)
}

// Text todo
type Text struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        string `xml:"MsgId"`
}

// PersonInfo todo
type PersonInfo struct {
	OpenID string

	StudentID string
	Name      string

	PhoneNumber string
}

//Get todo
func (p *PersonInfo) Get(openid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("PersonInfo"))
		data := bx.Get([]byte(openid))
		if data == nil || string(data) == "" {
			return nil
		}
		return json.Unmarshal(data, p)
	})

}

// Put todo
func (p *PersonInfo) Put() error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("PersonInfo"))
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return bx.Put([]byte(p.OpenID), data)
	})
}
