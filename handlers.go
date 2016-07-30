package main

import (
	"errors"
	"fmt"
	"strings"
)

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
	*Contact
}

// ServeMessage todo
func (c *contactHandler) ServeMessage(ctx *Context) {
	content := ctx.Message.Content
	if !strings.HasPrefix(content, "手机 ") {
		ctx.Next()
		return
	}

	name := content[len("手机 "):]
	err := c.Get(name)
	if err != nil {
		ctx.Errorln(err)
		ctx.Printf("没有%s的号码", name)
	}

	ctx.Printf("%s %s", name, c.PhoneNumber)
}

func handleBindPhone(ctx *Context) {
	content := ctx.Message.Content
	if !strings.HasPrefix(content, "我的手机 ") {
		ctx.Next()
		return
	}

	result := strings.Fields(content)
	name, phone := result[1], result[2]

	contact := Contact{
		PhoneNumber: phone,
	}

	err := contact.Put(name)
	if err != nil {
		ctx.Errorln(err)
		return
	}

	ctx.Printf("设置成功")

}
