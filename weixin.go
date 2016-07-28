package main

import "github.com/boltdb/bolt"

var db *bolt.DB

func main() {
	//设置数据库
	var err error
	db, err = bolt.Open("/root/data.db", 0600, nil)
	if err != nil {
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		buckets := []string{"default", "Contact", "Phone", "Person", "NameID"}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})

	//启动服务
	s := New()
	s.UseFunc(logHandler)
	s.UseFunc(testHandler)
	s.Run(":80")
}
