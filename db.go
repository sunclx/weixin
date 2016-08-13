package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func dbEditor() http.Handler {
	r := gin.Default()
	c := control{db}
	r.GET("/db", c.Index)
	r.GET("/buckets", c.Buckets)
	r.POST("/createBucket", c.CreateBucket)
	r.POST("/put", c.Put)
	r.POST("/get", c.Get)
	r.POST("/deleteKey", c.DeleteKey)
	r.POST("/deleteBucket", c.DeleteBucket)
	r.POST("/prefixScan", c.PrefixScan)
	r.GET("/web", sF.ServeFiles)
	return r
}

type control struct {
	db *bolt.DB
}

func (ctr *control) Index(c *gin.Context) {
	c.Redirect(301, "/web/html/layout.html")
}

func (ctr *control) CreateBucket(c *gin.Context) {
	if c.PostForm("bucket") == "" {
		c.String(200, "no bucket name | n")
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		return err
	})
	c.String(200, "ok")
}

func (ctr *control) DeleteBucket(c *gin.Context) {
	if c.PostForm("bucket") == "" {
		c.String(200, "no bucket name | n")
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(c.PostForm("bucket")))
		if err != nil {
			c.String(200, "error no such bucket | n")
		}
		return err
	})
	c.String(200, "ok")
}

func (ctr *control) DeleteKey(c *gin.Context) {
	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		c.String(200, "no bucket name or key | n")
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		if err != nil {
			c.String(200, "error no such bucket | n")
			return fmt.Errorf("bucket: %s", err)
		}

		err = b.Delete([]byte(c.PostForm("key")))
		if err != nil {
			c.String(200, "error Deleting KV | n")
			return fmt.Errorf("delete kv: %s", err)
		}

		return nil
	})
	c.String(200, "ok")
}

func (ctr *control) Put(c *gin.Context) {
	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		c.String(200, "no bucket name or key | n")
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.PostForm("bucket")))
		if err != nil {
			c.String(200, "error  creating bucket | n")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(c.PostForm("key")), []byte(c.PostForm("value")))
		if err != nil {
			c.String(200, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})
	c.String(200, "ok")
}

func (ctr *control) Get(c *gin.Context) {
	res := []string{"nok", ""}
	if c.PostForm("bucket") == "" || c.PostForm("key") == "" {
		res[1] = "no bucket name or key | n"
		c.JSON(200, res)
	}
	ctr.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.PostForm("bucket")))
		if b != nil {
			v := b.Get([]byte(c.PostForm("key")))
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

func (ctr *control) PrefixScan(c *gin.Context) {
	res := struct {
		Result string
		M      map[string]string
	}{Result: "nok"}
	res.M = make(map[string]string)
	if c.PostForm("bucket") == "" {
		res.Result = "no bucket name | n"
		c.JSON(200, res)
	}
	count := 0
	if c.PostForm("key") == "" {
		ctr.db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(c.PostForm("bucket")))
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
		ctr.db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(c.PostForm("bucket"))).Cursor()
			if b != nil {
				prefix := []byte(c.PostForm("key"))
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

func (ctr *control) Buckets(c *gin.Context) {
	res := []string{}
	ctr.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			b := []string{string(name)}
			res = append(res, b...)
			return nil
		})
	})
	c.JSON(200, res)
}
