package main

import (
	"errors"
	"fmt"

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

// Encode todo
func (n *Contact) Encode() []byte {
	return []byte(fmt.Sprintf("%s", n.PhoneNumber))
}

// Decode todo
func (n *Contact) Decode(data []byte) error {
	if n == nil {
		n = &Contact{}
	}
	_, err := fmt.Sscanf(string(data), "%s", &(n.PhoneNumber))
	return err
}

//Get todo
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

// Put todo
func (n *Contact) Put(openid string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("phone"))
		return bx.Put([]byte(openid), n.Encode())
	})
}
