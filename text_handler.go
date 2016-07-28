package main

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

type MessageHandler interface {
	ServeMessage(msg *Text)
}

const (
	PrefixPhone         string = "手机 "
	PrefixBindPhone     string = "我的手机 "
	PrefixBindStudentID string = "我的学号 "
	PrefixStudentID     string = "学号 "
)

func handlePhone(t *Text) string {
	if !strings.HasPrefix(t.Content, PrefixPhone) {
		return "wrong"
	}

	name := t.Content[len(PrefixPhone):]

	var n Contact
	err := n.Get(name)
	if err != nil {
		return fmt.Sprintf("没有%s的号码", name)
	}

	return fmt.Sprintf("%s %s", name, n.PhoneNumber)
}

func handleBindPhone(t *Text) string {
	content := t.Content
	result := strings.Fields(content)
	name, phone := result[1], result[2]

	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		err := bx.Put([]byte(name), []byte(phone))

		return err
	})

	return "设置成功"
}
