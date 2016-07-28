package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Server struct {
	appID     string
	token     string
	secruteID string

	handlers []Handler
	//	errorHandler ErrorHandler
}

func New() *Server {
	return &Server{
		appID:     cfg.AppID,
		token:     cfg.Token,
		secruteID: cfg.SecruteID,

		handlers: make([]Handler, 0, 8),
	}

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

	buffer := bytes.NewBuffer(nil)
	buffer.Reset()
	buffer.ReadFrom(r.Body)

	var t Text
	xml.Unmarshal(buffer.Bytes(), &t)

	// Context
	ctx := &Context{
		ResponseWriter: w,
		Request:        r,

		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		OpenID:    openid,

		Message: &t,

		index:    -1,
		handlers: s.handlers,
	}
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
