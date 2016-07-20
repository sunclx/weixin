package main

import "github.com/boltdb/bolt"

func db() {
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		return
	}
	defer db.Close()

}
