package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// User todo
type User struct {
	OpenID string

	StudentID string
	Name      string

	PhoneNumber string
}

//Get todo
func (p *User) Get(openid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("User"))
		data := bx.Get([]byte(openid))
		if data == nil || string(data) == "" {
			return nil
		}
		return json.Unmarshal(data, p)
	})

}

// Put todo
func (p *User) Put() error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("User"))
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return bx.Put([]byte(p.OpenID), data)
	})
}
