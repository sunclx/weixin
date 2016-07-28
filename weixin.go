package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kataras/iris"
)

func main() {
	go dbedit()

	s := New()

	s.UseFunc(mainHandle)
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

func mainHandle(c *Context) {
	//logConect(c)

	// //排除非POST请求
	// if c.MethodString() != iris.MethodPost {
	// 	c.WriteString("404")
	// 	return
	// }

	// //建议域名是否真确
	// if hostname := c.HostString(); hostname != "weixin.chenlixin.net" {
	// 	c.Log("异常域名:", hostname)
	// 	c.Write("404")
	// 	return
	// }

	// //检验是否是微信服务器的请求
	// signature := c.URLParam("signature")
	// timestamp := c.URLParam("timestamp")
	// nonce := c.URLParam("nonce")
	// if !validateURL(signature, timestamp, nonce) {
	// 	c.Log("参数错误%s,%s,%s.\n", signature, timestamp, nonce)
	// 	c.Write("404")
	// 	return
	// }

	//处理请求数据
	//handlerMux(c)
	fmt.Println("输入数据：")
	//fmt.Println(c)
	io.Copy(os.Stdout, c.Request.Body)

	template := `<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`

	s := fmt.Sprintf(template, c.OpenID, "gh_3fb3b0b8f2fa", time.Now().Unix(), "success")
	fmt.Println("输出数据：")
	fmt.Println(s)

	c.WriteString(s)
}

var verbose = false

// func updateHandle(c *iris.Context) {
// 	logConect(c)

// 	//执行命令
// 	cmd := exec.Command("git", "pull")
// 	cmd.Dir = "/root/go/src/github.com/sunclx/weixin"
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println("errors:", err)

// 	}
// 	c.WriteString("success")
// }

func logConect(c *iris.Context) {
	if !verbose {
		return
	}

	//记录请求
	fmt.Println("记录信息：", c.MethodString(), c.URI(), c.RemoteAddr())
}
