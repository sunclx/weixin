package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
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

func isValidateRequest(r *http.Request) bool {
	// 检查域名及请求方法
	if hostname := r.Host; r.Method != "POST" || hostname != "weixin.chenlixin.net" || r.URL.Path != "/" {
		fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
		return false
	}

	// 检验请求参数
	r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	if !validateURL(signature, timestamp, nonce, cfg.Token) {
		return false
	}
	return true
}
