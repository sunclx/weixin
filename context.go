package main

import "net/http"

type Handler interface {
	ServeMessage(ctx *Context)
}

//HandlerFunc ...
type HandlerFunc func(*Context)

var _ Handler = HandlerFunc(nil)

func (fn HandlerFunc) ServeMessage(ctx *Context) { fn(ctx) }

//Handlers ...
type Handlers struct {
	index    int
	handlers []Handler
}

func (h *Handlers) First(ctx *Context) {
	h.index = 0
	h.handlers[h.index].ServeMessage(ctx)
}
func (h *Handlers) Begin(ctx *Context) {
	h.handlers[h.index].ServeMessage(ctx)
}
func (h *Handlers) Next(ctx *Context) {
	h.index++
	if h.index >= len(h.handlers) {
		h.index--
		return
	}
	h.handlers[h.index].ServeMessage(ctx)
}

func (h *Handlers) IsEnd() bool {
	if h.index+1 >= len(h.handlers) {
		return true
	}
	return false
}

//Context ...
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	//QueryParams url.Values // 回调请求 URL 的查询参数集合
	//EncryptType  string     // 回调请求 URL 的加密方式参数: encrypt_type
	//MsgSignature string // 回调请求 URL 的消息体签名参数: msg_signature
	Signature string // 回调请求 URL 的签名参数: signature
	Timestamp string // 回调请求 URL 的时间戳参数: timestamp
	Nonce     string // 回调请求 URL 的随机数参数: nonce
	OpenID    string

	//MsgCiphertext []byte // 消息的密文文本
	//MsgPlaintext  []byte    // 消息的明文文本, xml格式
	//MixedMsg *MixedMsg // 消息

	//Token string // 当前消息所属公众号的 Token
	//AESKey      []byte // 当前消息加密所用的 aes-key, read-only!!!
	//Random      []byte // 当前消息加密所用的 random, 16-bytes
	//AppId       string // 当前消息加密所用的 AppId

	handlers *Handlers

	kvs map[string]interface{}
}

func (ctx *Context) Next() {
	ctx.handlers.Begin(ctx)
	for ctx.handlers.IsEnd() {
		ctx.handlers.Next(ctx)
	}
}
func (ctx *Context) Write(data []byte) (int, error) {
	return ctx.ResponseWriter.Write(data)
}

func (ctx *Context) WriteString(s string) {
	ctx.Write([]byte(s))
}
