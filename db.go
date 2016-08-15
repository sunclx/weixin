package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

func dbEditor() http.Handler {
	c := control{db, files}
	s := http.NewServeMux()
	s.HandleFunc("/db", c.Index)
	s.HandleFunc("/buckets", c.Buckets)
	s.HandleFunc("/createBucket", c.CreateBucket)
	s.HandleFunc("/put", c.Put)
	s.HandleFunc("/get", c.Get)
	s.HandleFunc("/deleteKey", c.DeleteKey)
	s.HandleFunc("/deleteBucket", c.DeleteBucket)
	s.HandleFunc("/prefixScan", c.PrefixScan)
	s.HandleFunc("/web/", c.Files)
	return s
}

type control struct {
	db    *bolt.DB
	files map[string]*staticFilesFile
}

func rstring(w http.ResponseWriter, content string) {
	w.WriteHeader(200)
	io.WriteString(w, content)
}

func rjson(w http.ResponseWriter, obj interface{}) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(obj)
}

func (ctr *control) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/web/html/layout.html", 301)
}

func (ctr *control) CreateBucket(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bucket := r.Form.Get("bucket")
	if bucket == "" {
		rstring(w, "no bucket name | n")
		return
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	rstring(w, "ok")
}

func (ctr *control) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bucket := r.Form.Get("bucket")
	if bucket == "" {
		rstring(w, "no bucket name | n")
		return
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		if err != nil {
			rstring(w, "error no such bucket | n")
		}
		return err
	})
	rstring(w, "ok")
}

func (ctr *control) DeleteKey(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bucket := r.Form.Get("bucket")
	key := r.Form.Get("key")
	if bucket == "" || key == "" {
		rstring(w, "no bucket name or key | n")
		return
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			rstring(w, "error no such bucket | n")
			return fmt.Errorf("bucket: %s", err)
		}

		err = b.Delete([]byte(key))
		if err != nil {
			rstring(w, "error Deleting KV | n")
			return fmt.Errorf("delete kv: %s", err)
		}

		return nil
	})
	rstring(w, "ok")
}

func (ctr *control) Put(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bucket := r.Form.Get("bucket")
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	if bucket == "" || key == "" {
		rstring(w, "no bucket name or key | n")
		return
	}
	ctr.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			rstring(w, "error  creating bucket | n")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(key), []byte(value))
		if err != nil {
			rstring(w, "error writing KV | n")
			return fmt.Errorf("create kv: %s", err)
		}

		return nil
	})
	rstring(w, "ok")
}

func (ctr *control) Get(w http.ResponseWriter, r *http.Request) {
	res := []string{"nok", ""}
	r.ParseForm()
	bucket := r.Form.Get("bucket")
	key := r.Form.Get("key")
	if bucket == "" || key == "" {
		res[1] = "no bucket name or key | n"
		rjson(w, res)
		return
	}
	ctr.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			v := b.Get([]byte(key))
			res[0] = "ok"
			res[1] = string(v)
			fmt.Printf("Key: %s\n", v)
		} else {
			res[1] = "error opening bucket| does it exist? | n"
		}
		return nil
	})
	rjson(w, res)
}

func (ctr *control) PrefixScan(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Result string
		M      map[string]string
	}{Result: "nok"}
	res.M = make(map[string]string)

	r.ParseForm()
	bucket := r.Form.Get("bucket")
	key := r.Form.Get("key")

	if bucket == "" {
		res.Result = "no bucket name | n"
		rjson(w, res)
		return
	}
	count := 0
	if key == "" {
		ctr.db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(bucket))
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
			b := tx.Bucket([]byte(bucket)).Cursor()
			if b != nil {
				prefix := []byte(key)
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
	rjson(w, res)
}

func (ctr *control) Buckets(w http.ResponseWriter, r *http.Request) {
	res := []string{}
	ctr.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			b := []string{string(name)}
			res = append(res, b...)
			return nil
		})
	})
	rjson(w, res)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

type staticFilesFile struct {
	data  string
	mime  string
	mtime time.Time
	// size is the size before compression. If 0, it means the data is uncompressed
	size int
	// hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
	hash string
}

func (ctr *control) Files(rw http.ResponseWriter, req *http.Request) {
	filename := strings.TrimPrefix(req.URL.Path, "/")
	if filename == "web/js//jquery-2.2.3.min.js" {
		filename = "web/js/jquery-2.2.3.min.js"
	}
	f, ok := ctr.files[filename]
	if !ok {
		http.NotFound(rw, req)
		return
	}
	header := rw.Header()
	if f.hash != "" {
		if hash := req.Header.Get("If-None-Match"); hash == f.hash {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("ETag", f.hash)
	}
	if !f.mtime.IsZero() {
		if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && f.mtime.Before(t.Add(1*time.Second)) {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("Last-Modified", f.mtime.UTC().Format(http.TimeFormat))
	}
	header.Set("Content-Type", f.mime)

	// Check if the asset is compressed in the binary
	if f.size == 0 {
		header.Set("Content-Length", strconv.Itoa(len(f.data)))
		io.WriteString(rw, f.data)
	} else {
		if header.Get("Content-Encoding") == "" && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			header.Set("Content-Encoding", "gzip")
			header.Set("Content-Length", strconv.Itoa(len(f.data)))
			io.WriteString(rw, f.data)
		} else {
			header.Set("Content-Length", strconv.Itoa(f.size))
			reader, _ := gzip.NewReader(strings.NewReader(f.data))
			io.Copy(rw, reader)
			reader.Close()
		}
	}
}
