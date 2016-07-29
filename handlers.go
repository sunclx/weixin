package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

// NameID todo
type NameID struct {
	Name      string
	StudentID string
}

// Encode todo
func (n *NameID) Encode() []byte {
	return []byte(fmt.Sprintf("%s&&%s", n.Name, n.StudentID))
}

// Decode todo
func (n *NameID) Decode(data []byte) error {
	if n == nil {
		n = &NameID{}
	}
	_, err := fmt.Sscanf(string(data), "%s&&%s", &(n.Name), &(n.StudentID))
	return err
}

// Get todo
func (n *NameID) Get(openid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("NameID"))
		data := bx.Get([]byte(openid))
		if data == nil {
			return errors.New("the openid doesn't exist")
		}
		return n.Decode(data)

	})
}

// Put todo
func (n *NameID) Put(openid string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("NameID"))
		return bx.Put([]byte(openid), n.Encode())
	})
}

// Contact todo
type Contact struct {
	PhoneNumber string
}

// ServeMessage todo
func (c *Contact) ServeMessage(ctx *Context) {
	content := ctx.Message.Content
	if !strings.HasPrefix(content, "手机 ") {
		ctx.Next()
		return
	}

	name := content[len("手机 "):]

	var n Contact
	err := n.Get(name)
	if err != nil {
		ctx.Printf("没有%s的号码", name)
	}

	ctx.Printf("%s %s", name, n.PhoneNumber)
}

// Encode todo
func (c *Contact) Encode() []byte {
	return []byte(fmt.Sprintf("%s", c.PhoneNumber))
}

// Decode todo
func (c *Contact) Decode(data []byte) error {
	if c == nil {
		c = &Contact{}
	}
	_, err := fmt.Sscanf(string(data), "%s", &(c.PhoneNumber))
	return err
}

//Get todo
func (c *Contact) Get(openid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		data := bx.Get([]byte(openid))
		if data == nil {
			return errors.New("the openid doesn't exist")
		}
		return c.Decode(data)

	})
}

// Put todo
func (c *Contact) Put(openid string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		return bx.Put([]byte(openid), c.Encode())
	})
}

func handleBindPhone(c *Context) {
	content := c.Message.Content
	if !strings.HasPrefix(content, "我的手机 ") {
		c.Next()
		return
	}

	result := strings.Fields(content)
	name, phone := result[1], result[2]

	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		err := bx.Put([]byte(name), []byte(phone))

		return err
	})

	c.Printf("设置成功")

}
