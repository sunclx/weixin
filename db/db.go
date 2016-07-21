package db

import (
	"sync"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("data.db", 0600, nil)
	if err != nil {
		panic(err)
	}

}

var pool = &sync.Pool{
	New: func() interface{} {

		return nil
	},
}

// func(text Text) (s string) {

// 		db, err := bolt.Open("data.db", 0600, nil)
// 		if err != nil {
// 			return ""
// 		}
// 		defer db.Close()
// 		content := text.Content
// 		if strings.HasPrefix(content, "我的学号是") {
// 			content = content[len(content)-8:]
// 			db.Update(func(tx *bolt.Tx) error {
// 				b := tx.Bucket([]byte("default"))
// 				err := b.Put([]byte(text.FromUserName), []byte(content))
// 				s = fmt.Sprintf("你的学号是%s，你是%s", content, "我们班的同学")

// 				return err
// 			})
// 		}

// 		db.Update(func(tx *bolt.Tx) error {
// 			b := tx.Bucket([]byte("default"))
// 			data := b.Get([]byte(text.FromUserName))
// 			if data == nil {
// 				s = `请输入"我的学号是00000000"`
// 				return nil
// 			}

// 			if string(data) == "09170515" {
// 				s = "你是跳跳，一个大美女"
// 				return nil
// 			}
// 			if string(data) == "09170512" {
// 				s = "你是乐乐，一个大美女"
// 				return nil
// 			}
// 			s = fmt.Sprintf("你的学号是%s，你是%s", data, "我们班的同学")

// 			return nil
// 		})

// 		return s
// 	}(text)

// 	return []byte(marshaMsg(text.FromUserName, resMsg))
// 	}
