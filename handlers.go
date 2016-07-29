package main

import (
	"strings"

	"github.com/boltdb/bolt"
)

func handlePhone(c *Context) {
	content := c.Message.Content
	if !strings.HasPrefix(content, "手机 ") {
		c.Next()
		return
	}

	name := content[len("手机 "):]

	var n Contact
	err := n.Get(name)
	if err != nil {
		c.Printf("没有%s的号码", name)
	}

	c.Printf("%s %s", name, n.PhoneNumber)
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
