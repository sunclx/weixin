package main

import (
	"fmt"
	"os/exec"

	"github.com/kataras/iris"
)

func main() {
	server := iris.New()

	//监听微信服务器的信息
	server.Post("/", func(c *iris.Context) {
		//记录请求
		fmt.Println(c.MethodString(), c.URI(), c.RemoteAddr())

		//建议域名是否真确
		if hostname := c.HostString(); hostname != "weixin.chenlixin.net" {
			c.Log("异常域名:", hostname)
			c.Write("404")
			return
		}

		//检验是否是微信服务器的请求
		signature := c.URLParam("signature")
		timestamp := c.URLParam("timestamp")
		nonce := c.URLParam("nonce")
		//openid:=c.URLParam("openid")
		if !validateURL(signature, timestamp, nonce) {
			c.Log("参数错误%s,%s,%s.\n", signature, timestamp, nonce)
			c.Write("404")
			return
		}

		//处理请求数据

		var t msgType
		c.ReadXML(&t)
		switch t.MsgType {
		case MsgTypeText:
			data := c.PostBody()
			c.Log("%s\n", data)
			TextHandle(data)
		default:
			c.Log("不支持该类型，%s.\n", t.MsgType)

		}
		fmt.Println("结束")
		c.Write("")

	})

	//监听github.com的自动更新
	server.HandleFunc("", "/update", func(c *iris.Context) {
		//记录请求
		fmt.Println(c.MethodString(), c.URI(), c.RemoteAddr())

		//执行命令
		cmd := exec.Command("go", "get", "-u", "github.com/sunclx/weixin")
		err := cmd.Run()
		if err != nil {
			fmt.Println("errors:", err)

		}
	})

	//启动服务
	server.Listen(":80")
}
