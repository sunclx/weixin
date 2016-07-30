package main

import (
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
		}
		name := parts[1]
		err := c.Get(name)
		if err != nil {
			ctx.Errorln(err)
			ctx.Printf("没有%s的号码%d", name, len(parts))
		}

		ctx.Printf("%s %s", name, c.PhoneNumber)
	case "我的手机":
		if len(parts) != 3 {
			ctx.Printf("参数错误")
		}

		name, phone := parts[1], parts[2]
		c.PhoneNumber = phone
		c.Put(name)

		ctx.Printf("设置成功")
	default:
		ctx.Next()
	}

}
