package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kataras/iris"
)

func main() {
	//go dbedit()

	s := New()
	s.UseFunc(logHandler)

	s.UseFunc(mainHandler)
	//server := iris.New()

	//监听github.com的自动更新
	//server.Post("/update", updateHandle)
	//server.Get("/update", updateHandle)

	//监听微信服务器的信息
	//server.HandleFunc("", "/", mainHandle)

	//启动服务
	//server.Listen(":80")
	s.Run(":80")
}

func logHandler(c *Context) {
	r := c.Request
	fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
}

func mainHandler(c *Context) {
	fmt.Println("输入数据是：")
	//fmt.Println(c)
	io.Copy(os.Stdout, c.Request.Body)

	s := fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`, c.OpenID, "gh_3fb3b0b8f2fa", time.Now().Unix(), "success")
	fmt.Println("输出数据是：")
	fmt.Println(s)

	c.WriteString(s)
}

var verbose = false

func logConect(c *iris.Context) {
	if !verbose {
		return
	}

	//记录请求
	fmt.Println("记录信息：", c.MethodString(), c.URI(), c.RemoteAddr())
}
