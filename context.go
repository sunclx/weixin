package main

import "net/http"

//Handler todo
type Handler interface {
	ServeMessage(ctx *Context)
}

//HandlerFunc ...
type HandlerFunc func(*Context)

var _ Handler = HandlerFunc(nil)

// ServeMessage todo
func (fn HandlerFunc) ServeMessage(ctx *Context) { fn(ctx) }

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

	index    int
	handlers []Handler

	//kvs map[string]interface{}
}

// Start todo
func (ctx *Context) Start() {
	ctx.handlers[ctx.index].ServeMessage(ctx)
}

// Next todo
func (ctx *Context) Next() {
	for ; ctx.index < len(ctx.handlers); ctx.index++ {
		ctx.handlers[ctx.index].ServeMessage(ctx)
	}
	ctx.index--
}

func (ctx *Context) Write(data []byte) (int, error) {
	return ctx.ResponseWriter.Write(data)
}

// WriteString todo
func (ctx *Context) WriteString(s string) {
	ctx.Write([]byte(s))
}
