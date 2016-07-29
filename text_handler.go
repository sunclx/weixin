package main

import (
	"strings"

	"github.com/boltdb/bolt"
)

func handlePhone(msg *Message) {
	t := msg.msg
	if !strings.HasPrefix(t.Content, "手机 ") {
		msg.Next()
		return
	}

	name := t.Content[len("手机 "):]

	var n Contact
	err := n.Get(name)
	if err != nil {
		msg.Printf("没有%s的号码", name)
	}

	msg.Printf("%s %s", name, n.PhoneNumber)
}

func handleBindPhone(msg *Message) {
	t := msg.msg
	if !strings.HasPrefix(t.Content, "我的手机 ") {
		msg.Next()
		return
	}

	content := t.Content
	result := strings.Fields(content)
	name, phone := result[1], result[2]

	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		err := bx.Put([]byte(name), []byte(phone))

		return err
	})

	msg.Printf("设置成功")

}
