package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type Server struct {
	appID     string
	token     string
	secruteID string

	handlers []Handler
	//	errorHandler ErrorHandler
}

func New() *Server {
	return &Server{token: "njmu0917"}

}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查域名及请求方法
	if hostname := r.Host; r.Method != "POST" || hostname != "weixin.chenlixin.net" || r.URL.Path != "/" {
		fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
		w.Write([]byte("404"))
		return
	}

	// 检验请求参数
	r.ParseForm()
	queryParams := r.Form
	signature := queryParams.Get("signature")
	timestamp := queryParams.Get("timestamp")
	nonce := queryParams.Get("nonce")
	openid := queryParams.Get("openid")
	if !validateURL(signature, timestamp, nonce, s.token) {
		w.Write([]byte("404"))
		return
	}

	// Context
	ctx := &Context{
		ResponseWriter: w,
		Request:        r,

		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		OpenID:    openid,

		handlers: make([]Handler, 0, 8),
	}
	ctx.handlers = s.handlers
	ctx.Next()
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) Use(handlers ...Handler) *Server {
	s.handlers = append(s.handlers, handlers...)
	return s
}

func (s *Server) UseFunc(handlersFunc ...func(ctx *Context)) *Server {

	for _, handlerFunc := range handlersFunc {
		s.handlers = append(s.handlers, HandlerFunc(handlerFunc))
	}

	return s
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
