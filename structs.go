package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

type person struct {
	StudentID  string
	Name       string
	Birthday   time.Time
	BirthPlace string
	Location   string
}

func personByByte(data []byte) *person {
	var p person
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil
	}
	return &p
}

func (p *person) JSON() string {
	s, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(s)
}
func (p *person) Get() {
	if p.StudentID == "" {
		p = nil
		return
	}

	err := db.View(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("persons"))
		bp := bx.Get([]byte(p.StudentID))
		p = personByByte(bp)
		return nil
	})
	if err != nil {
		p = nil
		return
	}
}
func (p *person) Put() {
	if p.StudentID == "" {
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		bx := tx.Bucket([]byte("persons"))
		return bx.Put([]byte(p.StudentID), []byte(p.JSON()))
	})

}
