package main

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

const (
	PrefixPhone         string = "手机 "
	PrefixBindPhone     string = "我的手机 "
	PrefixBindStudentID string = "我的学号 "
	PrefixStudentID     string = "学号 "
)

func handlePhone(t Text) string {
	content := t.Content
	name := content[len(PrefixPhone):]
	var msg string

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("phone"))
		return err
	})

	db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		phone := bx.Get([]byte(name))
		if phone == nil {
			msg = fmt.Sprintf("没有的%s号码", name)
			return nil
		}
		msg = fmt.Sprintf("%s %s", name, string(phone))

		return nil
	})
	return msg
}

func handleBindPhone(t Text) string {
	content := t.Content
	result := strings.Fields(content)
	name, phone := result[1], result[2]
	var msg string

	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		err := bx.Put([]byte(name), []byte(phone))

		return err
	})
	return msg
}

// db.Update(func(tx *bolt.Tx) error {
// 	content := msg.Content
// 	if strings.HasPrefix(content, "我的学号是") {
// 		content = content[len(content)-8:]
// 		b := tx.Bucket([]byte("default"))
// 		err := b.Put([]byte(msg.FromUserName), []byte(content))
// 		rmsg.Content = fmt.Sprintf("你的学号是%s\n", content)

// 		return err
// 	}

// 	b := tx.Bucket([]byte("default"))
// 	data := b.Get([]byte(msg.FromUserName))
// 	if data == nil {
// 		rmsg.Content = `请输入"我的学号是00000000"`
// 		return nil
// 	}

// 	if string(data) == "09170515" {
// 		rmsg.Content = "你是跳跳，一个大美女"
// 		return nil
// 	}
// 	if string(data) == "09170512" {
// 		rmsg.Content = "你是乐乐，一个大美女"
// 		return nil
// 	}
// 	rmsg.Content = fmt.Sprintf("你的学号是%s，你是%s", data, "我们班的同学")

// 	return nil
// })

// c.WriteString(rmsg.String())
