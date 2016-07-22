package main

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

type NameID struct {
	Name      string
	StudentID string
}

func (n *NameID) Encode() []byte {
	return []byte(fmt.Sprintf("%s&&%s", n.Name, n.StudentID))
}

func (n *NameID) Decode(data []byte) error {
	if n == nil {
		n = &NameID{}
	}
	_, err := fmt.Sscanf(string(data), "%s&&%s", &(n.Name), &(n.StudentID))
	return err
}

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

func (n *NameID) Put(openid string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("NameID"))
		return bx.Put([]byte(openid), n.Encode())
	})
}

type Contact struct {
	PhoneNumber string
}

func (n *Contact) Encode() []byte {
	return []byte(fmt.Sprintf("%s", n.PhoneNumber))
}

func (n *Contact) Decode(data []byte) error {
	if n == nil {
		n = &Contact{}
	}
	_, err := fmt.Sscanf(string(data), "%s", &(n.PhoneNumber))
	return err
}

func (n *Contact) Get(openid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		data := bx.Get([]byte(openid))
		if data == nil {
			return errors.New("the openid doesn't exist")
		}
		return n.Decode(data)

	})
}

func (n *Contact) Put(openid string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		return bx.Put([]byte(openid), n.Encode())
	})
}
