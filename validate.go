package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	token string = "njmu0917"
)

func validate(r *http.Request) bool {
	if hostname := r.URL.Host; hostname != "weixin.chenlixin.net" {
		return false
	}

	//获取参数
	values := r.URL.Query()
	signature := values.Get("signature")
	timestamp := values.Get("timestamp")
	nonce := values.Get("nonce")

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
