package main

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	"github.com/kataras/iris"
)

const (
	token string = "njmu0917"
)

func validateURL(parameters iris.PathParameters) bool {

	//获取参数

	signature := parameters.Get("signature")
	timestamp := parameters.Get("timestamp")
	nonce := parameters.Get("nonce")

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
