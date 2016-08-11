package main

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

func fatalError(err error) {
	if err != nil {
		lg.Fatalln(err)
	}
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
