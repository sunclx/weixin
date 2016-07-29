package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var cfg *config
var f io.Writer

type config struct {
	DeveloperID string
	AppID       string
	Token       string
	SecruteID   string
	DBPath      string
}

func init() {
	//初始化配置
	cfg = &config{
		DeveloperID: "gh_3fb3b0b8f2fa",
		AppID:       "",
		Token:       "njmu0917",
		SecruteID:   "",
		DBPath:      "/root/data.db",
	}

	//设置数据库
	db, _ = bolt.Open(cfg.DBPath, 0600, nil)

	f, _ = os.Create("/root/log.log")

}

//验证函数
func validateURL(signature, timestamp, nonce, token string) bool {
	//排序参数并合并
	ss := []string{token, timestamp, nonce}
	sort.Strings(ss)
	s := strings.Join(ss, "")

	//计算sha1的值
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	//比较计算的signature与获取值比较
	if signatureHex := fmt.Sprintf("%x", bs); signatureHex != signature {
		return false
	}
	return true
}
