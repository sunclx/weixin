package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Handler todo
type Handler interface {
	ServeMessage(c *Context)
}

// HandlerFunc todo
type HandlerFunc func(c *Context)

// ServeMessage todo
func (fn HandlerFunc) ServeMessage(c *Context) {
	fn(c)
}

// Contact todo
type Contact struct {
	PhoneNumber string
}

//Get todo
func (c *Contact) Get(openid string) error {
	bt, err := bx.New([]byte("Contacts"))
	if err != nil {
		return err
	}

	data, err := bt.Get([]byte(openid))
	if err != nil {
		return err
	}

	if c == nil {
		c = &Contact{}
	}
	_, err = fmt.Sscanf(string(data), "%s", &(c.PhoneNumber))
	return err

}

// Put todo
func (c *Contact) Put(openid string) error {
	if c == nil {
		return errors.New("empty contact")
	}

	bt, err := bx.New([]byte("Contacts"))
	if err != nil {
		return err
	}

	return bt.Put([]byte(openid), []byte(fmt.Sprintf("%s", c.PhoneNumber)))
}

type contactHandler struct {
	Contact
}

// ServeMessage todo
func (c *contactHandler) ServeMessage(ctx *Context) {
	content := ctx.Message.Content
	parts := strings.Fields(content)

	switch parts[0] {
	case "手机":
		if len(parts) != 2 {
			ctx.Printf("参数错误")
			ctx.Infoln(parts)
			return
		}
		name := parts[1]
		err := c.Get(name)
		if err != nil {
			ctx.WithError(err).Info(parts)
			ctx.Printf("没有%s的号码", name)
			return
		}

		ctx.Printf("%s %s", name, c.PhoneNumber)
	case "我的手机":
		if len(parts) != 3 {
			ctx.Printf("参数错误")
			ctx.Infoln(parts)
			return
		}

		name, phone := parts[1], parts[2]
		c.PhoneNumber = phone
		c.Put(name)

		ctx.Printf("设置成功")
	default:
		ctx.Next()
	}

}

type NameID struct {
	Name      string
	StudentID string
}

func (n *NameID) Get(openid string) error {
	bt, err := bx.New([]byte("NameID"))
	if err != nil {
		return err
	}

	data, err := bt.Get([]byte(openid))
	if data == nil || err != nil {
		return errors.New("Not Exist")
	}

	parts := bytes.Split(data, []byte("&&"))
	if len(parts) != 2 {
		return errors.New("wrong data")
	}

	n.Name = string(parts[0])
	n.StudentID = string(parts[1])
	return nil
}
func (n *NameID) Put(openid string) error {
	bt, err := bx.New([]byte("NameID"))
	if err != nil {
		return err
	}

	return bt.Put([]byte(openid), []byte(n.Name+"&&"+n.StudentID))
}

type openidHandler struct {
	NameID
}

func (o *openidHandler) ServeMessage(ctx *Context) {
	content := ctx.Message.Content
	parts := strings.Fields(content)
	switch parts[0] {
	case "我的学号":
		if len(parts) != 2 {
			ctx.Printf("参数错误")
			return
		}
		o.Name = parts[1]
		o.StudentID = parts[1]
		o.Put(ctx.OpenID)
		ctx.Printf("学号绑定成功")
		return
	}

	err := o.Get(ctx.OpenID)
	if err != nil {
		ctx.Printf(`请输入“我的学号 00000000”`)
		return
	}

	ctx.Next()
}
