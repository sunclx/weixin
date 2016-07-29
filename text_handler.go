package main

import (
	"strings"

	"github.com/boltdb/bolt"
)

const (
	PrefixPhone         string = "手机 "
	PrefixBindPhone     string = "我的手机 "
	PrefixBindStudentID string = "我的学号 "
	PrefixStudentID     string = "学号 "
)

func handlePhone(msg *Message) {
	t := msg.msg
	if !strings.HasPrefix(t.Content, PrefixPhone) {
		msg.Next()
		return
	}

	name := t.Content[len(PrefixPhone):]

	var n Contact
	err := n.Get(name)
	if err != nil {
		msg.Printf("没有%s的号码", name)
	}

	msg.Printf("%s %s", name, n.PhoneNumber)
}

func handleBindPhone(msg *Message) {
	t := msg.msg
	if !strings.HasPrefix(t.Content, PrefixBindPhone) {
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
