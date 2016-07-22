package main

import "github.com/boltdb/bolt"

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("/root/ddata.db", 0600, nil)
	if err != nil {
		return
	}
}
