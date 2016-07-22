package main

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
)

func dbedit() {

	r := iris.New()

	r.Get("/", boltbrowserweb.Index)

	r.Get("/buckets", boltbrowserweb.Buckets)
	r.Post("/createBucket", boltbrowserweb.CreateBucket)
	r.Post("/put", boltbrowserweb.Put)
	r.Post("/get", boltbrowserweb.Get)
	r.Post("/deleteKey", boltbrowserweb.DeleteKey)
	r.Post("/deleteBucket", boltbrowserweb.DeleteBucket)
	r.Post("/prefixScan", boltbrowserweb.PrefixScan)

	r.Static("/web", "./", -1)

	r.Listen(":8080")

}

type boltBrowerWeb struct{}

var boltbrowserweb *boltBrowerWeb

func (b *boltBrowerWeb) Index(c *iris.Context) {
	logConect(c)
	c.Redirect("/web/html/layout.html", 301)

}

func (b *boltBrowerWeb) CreateBucket(c *iris.Context) {
	logConect(c)

	if c.FormValueString("bucket") == "" {
		c.Text(200, "no bucket name | n")

	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.FormValueString("bucket")))
		b = b
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	c.Text(200, "ok")

}

func (b *boltBrowerWeb) DeleteBucket(c *iris.Context) {
	logConect(c)

	if c.FormValueString("bucket") == "" {
		c.Text(200, "no bucket name | n")
	}

	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(c.FormValueString("bucket")))

		if err != nil {

			c.Text(200, "error no such bucket | n")
			return fmt.Errorf("bucket: %s", err)
		}

		return nil
	})

	c.Text(200, "ok")

}

func (b *boltBrowerWeb) DeleteKey(c *iris.Context) {
	logConect(c)

	if c.FormValueString("bucket") == "" || c.FormValueString("key") == "" {
		c.Text(200, "no bucket name or key | n")
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.FormValueString("bucket")))
		b = b
		if err != nil {

			c.Text(200, "error no such bucket | n")
			return fmt.Errorf("bucket: %s", err)
		}

		err = b.Delete([]byte(c.FormValueString("key")))

		if err != nil {

			c.Text(200, "error Deleting KV | n")
			return fmt.Errorf("delete kv: %s", err)
		}

		return nil
	})

	c.Text(200, "ok")

}

func (b *boltBrowerWeb) Put(c *iris.Context) {
	logConect(c)

	if c.FormValueString("bucket") == "" || c.FormValueString("key") == "" {
		c.Text(200, "no bucket name or key | n")
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.FormValueString("bucket")))
		b = b
		if err != nil {

			c.Text(200, "error  creating bucket | n")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(c.FormValueString("key")), []byte(c.FormValueString("value")))

		if err != nil {

			c.Text(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})

	c.Text(200, "ok")

}

func (b *boltBrowerWeb) Get(c *iris.Context) {
	logConect(c)

	res := []string{"nok", ""}

	if c.FormValueString("bucket") == "" || c.FormValueString("key") == "" {

		res[1] = "no bucket name or key | n"
		c.JSON(200, res)
	}

	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(c.FormValueString("bucket")))

		if b != nil {

			v := b.Get([]byte(c.FormValueString("key")))

			res[0] = "ok"
			res[1] = string(v)

			fmt.Printf("Key: %s\n", v)

		} else {

			res[1] = "error opening bucket| does it exist? | n"

		}
		return nil

	})

	c.JSON(200, res)

}

type Result struct {
	Result string
	M      map[string]string
}

func (b *boltBrowerWeb) PrefixScan(c *iris.Context) {
	logConect(c)

	res := Result{Result: "nok"}

	res.M = make(map[string]string)

	if c.FormValueString("bucket") == "" {

		res.Result = "no bucket name | n"
		c.JSON(200, res)
	}

	count := 0

	if c.FormValueString("key") == "" {

		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(c.FormValueString("bucket")))

			if b != nil {

				c := b.Cursor()

				for k, v := c.First(); k != nil; k, v = c.Next() {
					res.M[string(k)] = string(v)

					if count > 2000 {
						break
					}
					count++
				}

				res.Result = "ok"

			} else {

				res.Result = "no such bucket available | n"

			}

			return nil
		})

	} else {

		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(c.FormValueString("bucket"))).Cursor()

			if b != nil {

				prefix := []byte(c.FormValueString("key"))

				for k, v := b.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = b.Next() {
					res.M[string(k)] = string(v)
					if count > 2000 {
						break
					}
					count++
				}

				res.Result = "ok"

			} else {

				res.Result = "no such bucket available | n"

			}

			return nil
		})

	}

	c.JSON(200, res)

}

func (b *boltBrowerWeb) Buckets(c *iris.Context) {
	logConect(c)

	res := []string{}

	db.View(func(tx *bolt.Tx) error {

		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {

			b := []string{string(name)}
			res = append(res, b...)
			return nil
		})

	})

	c.JSON(200, res)

}
