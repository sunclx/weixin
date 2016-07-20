package main

import "github.com/boltdb/bolt"

func db() {
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		return
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("default"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("oXU74wMt--mr4eaVyhmJ5h_lSJP0"), []byte("09170510"))
		return err
	})

}
